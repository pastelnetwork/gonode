import argparse
import math
import shutil

from tqdm import tqdm
import cv2
import torch
import random
import numpy as np
from sklearn.manifold import TSNE
import matplotlib.pyplot as plt

from DataSetClasses import OriginalImagesDataset, collate_skip_empty, colors_per_class
from ImageWithCoordinates import ImageWithCoordinates, EuclideanDistanceMap
from resnet import ResNet101

#GLobal varaibles / switches
euclidean_distance_map = None

TO_BE_DELETED_PATH= 'dataset/to_be_deleted/'

# To get the same results for every run on the same images
def fix_random_seeds():
    seed = 10
    random.seed(seed)
    torch.manual_seed(seed)
    np.random.seed(seed)


def get_features(dataset, batch, num_images):
    # move the input and model to GPU for speed if available
    if torch.cuda.is_available():
        device = 'cuda'
    else:
        device = 'cpu'

    # initialize our implementation of ResNet
    model = ResNet101(pretrained=True)
    model.eval()
    model.to(device)

    # read the dataset and initialize the data loader
    dataset = OriginalImagesDataset(dataset, num_images)
    print("dataset:{}".format(dataset))
    print("images:{}".format(num_images))
    dataloader = torch.utils.data.DataLoader(dataset, batch_size=batch, collate_fn=collate_skip_empty, shuffle=True) #batch

    # we'll store the features as NumPy array of size num_images x feature_size
    features = None

    # we'll also store the image labels and paths to visualize them later
    labels = []
    image_paths = []

    for batch in tqdm(dataloader, desc='Running the model inference'):
        images = batch['image'].to(device)
        labels += batch['label']
        image_paths += batch['image_path']

        with torch.no_grad():
            output = model.forward(images)

        current_features = output.cpu().numpy()
        if features is not None:
            features = np.concatenate((features, current_features))
        else:
            features = current_features

    return features, labels, image_paths


# scale and move the coordinates so they fit [0; 1] range
def scale_to_01_range(x):
    # compute the distribution range
    value_range = (np.max(x) - np.min(x))

    # move the distribution so that it starts from zero
    # by extracting the minimal value from all its values
    starts_from_zero = x - np.min(x)

    # make the distribution fit [0; 1] by dividing by its range
    return starts_from_zero / value_range


def scale_image(image, max_image_size):
    image_height, image_width, _ = image.shape

    scale = max(1, image_width / max_image_size, image_height / max_image_size)
    image_width = int(image_width / scale)
    image_height = int(image_height / scale)

    image = cv2.resize(image, (image_width, image_height))
    return image


def draw_rectangle_by_class(image, label):
    image_height, image_width, _ = image.shape

    # get the color corresponding to image class
    color = colors_per_class[label]
    image = cv2.rectangle(image, (0, 0), (image_width - 1, image_height - 1), color=color, thickness=5)

    return image


def compute_plot_coordinates(image, x, y, image_centers_area_size, offset):
    image_height, image_width, _ = image.shape

    # compute the image center coordinates on the plot
    center_x = int(image_centers_area_size * x) + offset

    # in matplotlib, the y axis is directed upward
    # to have the same here, we need to mirror the y coordinate
    center_y = int(image_centers_area_size * (1 - y)) + offset

    # knowing the image center, compute the coordinates of the top left and bottom right corner
    tl_x = center_x - int(image_width / 2)
    tl_y = center_y - int(image_height / 2)

    br_x = tl_x + image_width
    br_y = tl_y + image_height

    return tl_x, tl_y, br_x, br_y


# Keep in mind that we do not have such feature in 'live'.
# So we won't know the image's parent (the original) because
# this is what we are trying to reach
def move_to_delete_folder(image_path):
    filename = image_path.rsplit('/',1)
    #for splitings in filename:
        #print (splitings)
    shutil.copyfile(image_path, TO_BE_DELETED_PATH+filename[-1])
    pass


def is_parents_euclidean_distance_big(child_img, images, tx, ty):
    # now we'll put a small copy of every image to its corresponding T-SNE coordinate

    for image_path, x, y in tqdm(
            zip(images, tx, ty),
            desc='Searching for parents',
            total=len(images)
    ):

        if image_path == child_img.parentName:
            distance = math.sqrt((x - child_img.x)**2+(y-child_img.y)**2)
            print('Euclidean distance between parent{} and child{}: {}'.format(child_img.parentName, child_img.imageName, distance))

            # Not reliable :'(
            if(distance > 0.5):
                print ("Shall be deleted")
                #move_to_delete_folder(child_img.imageName)
            return


# Building up the EuclideanDistanceMap class instance to be able to search via
# filenames
def build_up_euclidean_distance_matrix(list_of_image_coordinates):
    global euclidean_distance_map
    # build a complex array of your cells
    z = np.array([complex(image.x, image.y) for image in list_of_image_coordinates])

    # Euclidean distance numpy array
    # mesh this array so that you will have all combinations
    m, n = np.meshgrid(z, z)
    # get the distance via the norm
    euclidean_map = abs(m - n)

    # Init class euclidean_distance_map
    euclidean_distance_map = EuclideanDistanceMap(list_of_image_coordinates, euclidean_map)
    print('Euclidean average distance: {}'.format(euclidean_distance_map.get_mean_eu_distance()))


# Get the names of the images which are the top N closest in terms of tSNE
def retrieve_closest_N_images_per_image(image_name, number):
    list_filenames = []
    list_topN_closests_file_idx=[]
    if euclidean_distance_map:
        #Getting the index value we use the find the 'column' in the Euclidean distance array
        index = next((i for i, item in enumerate(euclidean_distance_map.image_list) if item.imageName == image_name), -1)

        if index != -1:
            temp_list = euclidean_distance_map.eu_distances[:,index]

            # Remove zeros because it is 'itself'
            # This list will contain all the indexes which are the top N closest to the requested image
            delted_itself = np.delete(temp_list, np.where(temp_list == [0]), axis=0)
            # Now we need the top 'number' indexes with which we can locate the filenames:
            list_topN_closests_file_idx = delted_itself.argsort()[:number]

            # Now that we have the indexes, let's collect the names
            for idx in list_topN_closests_file_idx:
                list_filenames.append(euclidean_distance_map.image_list[idx].imageName)

    print("These are the top {} images close to this {} image:\n".format(number,image_name))
    for image in list_filenames:
        print ("{}".format(image))



# Keep in mind that we do not have such feature in 'live'.
# So we won't know the image's parent (the original) because
# this is what we are trying to reach
def get_euclidean_distances_of_transformed_images(list_of_images_where_transformation_happenned):
    global euclidean_distance_map

    if euclidean_distance_map:
        for image in list_of_images_where_transformation_happenned:
            distance = euclidean_distance_map.get_eu_distance_bw_2_images(image.imageName, image.parentName)
            print ("Euclidean distance between {} and child : {} is : {}".format(image.parentName, image.imageName, distance))


def visualize_tsne_images(tx, ty, images, labels, plot_size=1000, max_image_size=100):
    # we'll put the image centers in the central area of the plot
    # and use offsets to make sure the images fit the plot
    offset = max_image_size // 2
    image_centers_area_size = plot_size - 2 * offset

    tsne_plot = 255 * np.ones((plot_size, plot_size, 3), np.uint8)
    list_of_image_coordinates = []
    list_of_images_where_transformation_happenned = []

    # now we'll put a small copy of every image to its corresponding T-SNE coordinate
    for image_path, label, x, y in tqdm(
            zip(images, labels, tx, ty),
            desc='Building the T-SNE plot',
            total=len(images)
    ):
        #print("How many images we have: {}".format(len(images)))
        image_w_coord = ImageWithCoordinates(image_path, x, y)
        list_of_image_coordinates.append(image_w_coord)
        # Search for parent and get the euclidean distance.
        # If big enough we can say that it shall be thrown.

        # Keep in mind that in live we don't have the possibility
        # to have the child-parent images defined. This is exactly
        # we are trying to reach, so labellig with names is just
        # for data sampling
        if not image_w_coord.isOriginal:
            # Debug print
            #print ("Child found")
            is_parents_euclidean_distance_big(image_w_coord, images, tx, ty)
            list_of_images_where_transformation_happenned.append(image_w_coord)

        image = cv2.imread(image_path)

        # scale the image to put it to the plot
        image = scale_image(image, max_image_size)

        # draw a rectangle with a color corresponding to the image class
        image = draw_rectangle_by_class(image, label)

        # compute the coordinates of the image on the scaled plot visualization
        tl_x, tl_y, br_x, br_y = compute_plot_coordinates(image, x, y, image_centers_area_size, offset)

        # put the image to its TSNE coordinates using numpy subarray indices
        tsne_plot[tl_y:br_y, tl_x:br_x, :] = image

    build_up_euclidean_distance_matrix(list_of_image_coordinates)
    # Example
    # retrieve_closest_N_images_per_image('dataset/opensea_dataset/openseaio_reduced/reduced/https___opensea.io_assets_0x7c40c393dc0f283f318791d746d894ddd3693572_8246.png', 15)

    # Keep in mind that we do not have such information 'live' as original vs. fake
    # For creating dataset, sure !
    get_euclidean_distances_of_transformed_images(list_of_images_where_transformation_happenned)

    plt.imshow(tsne_plot[:, :, ::-1])
    plt.show()


def visualize_tsne_points(tx, ty, labels):
    # initialize matplotlib plot
    fig = plt.figure()
    ax = fig.add_subplot(111)

    # for every class, we'll add a scatter plot separately
    for label in colors_per_class:
        # find the samples of the current class in the data
        indices = [i for i, l in enumerate(labels) if l == label]

        # extract the coordinates of the points of this class only
        current_tx = np.take(tx, indices)
        current_ty = np.take(ty, indices)

        # convert the class color to matplotlib format:
        # BGR -> RGB, divide by 255, convert to np.array
        color = np.array([colors_per_class[label][::-1]], dtype=float) / 255

        # add a scatter plot with the correponding color and label
        ax.scatter(current_tx, current_ty, c=color, label=label)

    # build a legend using the labels we set previously
    ax.legend(loc='best')

    # finally, show the plot
    plt.show()


def visualize_tsne(tsne, images, labels, plot_size=1000, max_image_size=100):
    # extract x and y coordinates representing the positions of the images on T-SNE plot
    tx = tsne[:, 0]
    ty = tsne[:, 1]

    # scale and move the coordinates so they fit [0; 1] range
    tx = scale_to_01_range(tx)
    ty = scale_to_01_range(ty)

    # visualize the plot: samples as colored points
    visualize_tsne_points(tx, ty, labels)

    # visualize the plot: samples as images
    visualize_tsne_images(tx, ty, images, labels, plot_size=plot_size, max_image_size=max_image_size)


def do_tsne_calc():
    parser = argparse.ArgumentParser()

    parser.add_argument('--path', type=str, default='dataset/trash_automation')
    parser.add_argument('--batch', type=int, default=64)
    parser.add_argument('--num_images', type=int, default=500)
    args = parser.parse_args()

    fix_random_seeds()

    features, labels, image_paths = get_features(
        dataset=args.path,
        batch=args.batch,
        num_images=args.num_images
    )

    tsne = TSNE(n_components=2).fit_transform(features)

    visualize_tsne(tsne, image_paths, labels)

if __name__ == '__main__':
    do_tsne_calc()