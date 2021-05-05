# import the necessary packages
import argparse
import os
from random import randrange
from PIL import Image, ImageOps, ImageFilter, ImageEnhance
import PIL
import numpy as np

# Anaconda package was not existing before and porting a package from Pypi
# was the easiest way to have lots of instagram filters in an easy-to-apply format.
# Package is available at conda (under my acct name): https://anaconda.org/adaki2004/pilgram
# Anyone (who uses conda environment (linux-64 currently) can install now pilgram
# by running the following command in terminal 'conda install -c adaki2004 pilgram'
import pilgram

# Global variables / switches
# If switch is 0 then we shall fill up
# the constants in def parse_folders()
USE_CMD_ARGS = 0

# Input - output folders
ORIGINAL_FILE_PATH = ''
TRANSFORMED_FILE_PATH = ''

# Crop size ratio used by random_crop function
CROP_RATIO = 0.3
# Stretch width ratio used by random_stretch function
STRETCH_WIDTH_RATIO = 0.4
# Stretch height ratio used by random_stretch function
STRETCH_HEIGHT_RATIO = 0.2
# Radius parameter for GaussianBlur
BLUR_RADIUS = 15
# Determines the scaling factor for step1 in pixelation
PIXELATE_SCALING = 0.1
# Enhancement ratio: (min 0, max 1)
# Factor 1.0 always returns a copy of the original image.
# Lower factors mean less color (brightness, contrast, etc).
ENHANCEMENT_FACTOR = 0.7
# Mosaic: resize the input to fit original image size?
RESIZE_INPUT = True
# Mosaic: grid size. THe bigger this number the the more detailed the image (and slower the process)
GRID_SIZE = 100
# Mosaic: re-use any image in input
REUSE_IMAGES = True
# Mosaic: List of images which will be contained in the mosaic
INPUT_IMG_LIST = []


def insta_filter_types(i):
    insta_filter_map = {
        0: '_1977',
        1: 'aden',
        2: 'brannan',
        3: 'brooklyn',
        4: 'clarendon',
        5: 'earlybird',
        6: 'gingham',
        7: 'hudson',
        8: 'inkwell',
        9: 'kelvin',
        10: 'lark',
        11: 'lofi',
        12: 'maven',
        13: 'mayfair',
        14: 'moon',
        15: 'nashville',
        16: 'perpetua',
        17: 'reyes',
        18: 'rise',
        19: 'slumber',
        20: 'stinson',
        21: 'toaster',
        22: 'valencia',
        23: 'walden',
        24: 'willow',
        25: 'xpro2'
    }
    return insta_filter_map.get(i, 'Invalid')


def base_filter_types(i):
    filter_map = {
        0: 'BLUR',
        1: 'CONTOUR',
        2: 'DETAIL',
        3: 'EDGE_ENHANCE',
        4: 'EDGE_ENHANCE_MORE',
        5: 'EMBOSS',
        6: 'SHARPEN',
        7: 'SMOOTH',
        8: 'SMOOTH_MORE'
    }
    return filter_map.get(i, "Invalid")


def enhancement_types(i):
    enhancement_map = {
        0: 'Color',
        1: 'Contrast',
        2: 'Brightness',
    }
    return enhancement_map.get(i, "Invalid")

def parse_folders():
    global ORIGINAL_FILE_PATH
    global TRANSFORMED_FILE_PATH

    if USE_CMD_ARGS:
        # construct the argument parse and parse the arguments
        ap = argparse.ArgumentParser()
        ap.add_argument("-i", "--inputFolder", type=str, required=True,
                        help="path to optional input folder")
        ap.add_argument("-o", "--outputFolder", type=str, required=True,
                        help="path to optional output folder")
        args = vars(ap.parse_args())

        ORIGINAL_FILE_PATH = args["inputFolder"]
        TRANSFORMED_FILE_PATH = args["outputFolder"]
    else:
        ORIGINAL_FILE_PATH = "input_folder/"
        TRANSFORMED_FILE_PATH = "output_folder/"


# Function which shall be used to save images
# based on transformation
def save_image(transform_type, filename, to_be_saved):
    file = os.path.splitext(filename)[0]
    file_extension = os.path.splitext(filename)[1]

    new_filename = file + transform_type + file_extension
    to_be_saved.save(TRANSFORMED_FILE_PATH + new_filename)
    to_be_saved.close()


def crop(img, filename):
    print("Cropping image")
    # Make a copy which will be saved later on
    cropped_img = img.copy()

    # get image size
    x, y = cropped_img.size

    # determine the mask for cropping
    # no i pick up the 30% of the width
    # and 30% of the height
    matrix_x = int(x * CROP_RATIO)
    matrix_y = int(x * CROP_RATIO)
    x1 = randrange(0, x - matrix_x)
    y1 = randrange(0, y - matrix_y)

    cropped_img = cropped_img.crop((x1, y1, x1 + matrix_x, y1 + matrix_y))

    # Save image
    save_image('_cropped', filename, cropped_img)


def flip(img, filename):
    print("Flipping image")
    # Make a copy which will be saved later on
    flipped_img = img.copy()

    # random nr. between horizontally vs. vertically flipping
    irand = randrange(0, 2)
    if irand == 0:
        # flip image vertically
        flipped_img = flipped_img.transpose(PIL.Image.FLIP_LEFT_RIGHT)
    else:
        # flip image horizontally
        flipped_img = flipped_img.transpose(PIL.Image.FLIP_TOP_BOTTOM)

    save_image('_flipped', filename, flipped_img)


def stretch(img, filename):
    print("Stretching image")
    # Make a copy which will be saved later on
    stretched_img = img.copy()

    # get image size
    x, y = stretched_img.size

    new_width = int(x * STRETCH_WIDTH_RATIO)
    new_height = int(y * STRETCH_HEIGHT_RATIO)

    stretched_img = stretched_img.resize((new_width, new_height), Image.ANTIALIAS)

    save_image('_stretched', filename, stretched_img)


def base_filter(img, filename):
    # Make a copy which will be saved later on
    filtered_img = img.copy()
    irand = randrange(0, 9)

    filter_type = base_filter_types(irand)
    filter_attribute: object = getattr(ImageFilter, filter_type)
    print("Image is filtered based on: {} base filter".format(filter_type))
    if filter_type == 'BLUR':
        filtered_img = filtered_img.filter(ImageFilter.GaussianBlur(BLUR_RADIUS))
    else:
        filtered_img = filtered_img.filter(filter_attribute)

    save_image('_base_filtered_{}'.format(filter_type), filename, filtered_img)


def insta_filter(img, filename):
    # Make a copy which will be saved later on
    filtered_img = img.copy()
    irand = randrange(0, 26)

    filter_type = insta_filter_types(irand)

    # Trying to create/access function from dictionary string
    try:
        func = getattr(pilgram, filter_type)
        filtered_img = func(filtered_img)
        print("Image is filtered based on: {} insta filter".format(filter_type))
    except AttributeError:
        print("Function not found")

    save_image('_insta_filtered_{}'.format(filter_type), filename, filtered_img)


def edges(img, filename):
    print("Finding edges")
    # Make a copy which will be saved later on
    filtered_img = img.copy()

    filtered_img = filtered_img.filter(ImageFilter.FIND_EDGES)

    save_image('_edges', filename, filtered_img)


def invert(img, filename):
    print("Invert")
    # Make a copy which will be saved later on
    inverted_img = img.copy()

    inverted_img = ImageOps.invert(inverted_img)

    save_image('_inverted', filename, inverted_img)


def image_enhance(img, filename):
    print("Saturate")
    # Make a copy which will be saved later on
    saturated_img = img.copy()

    irand = randrange(0, 3)

    filter_type = enhancement_types(irand)
    filter_attribute = getattr(PIL.ImageEnhance, filter_type)
    print("Image is filtered based on: {} base filter".format(filter_type))

    converter = filter_attribute(saturated_img)
    saturated_img = converter.enhance(ENHANCEMENT_FACTOR)

    save_image('_enhanced_{}'.format(filter_type), filename, saturated_img)


def pixelate(img, filename):
    # Simply put pixelation is:
    # 1. Resize to a smaller size with the help of PIXELATE_SCALING (smaller the number the more pixelated the result
    # 2. Scale it back to original using nearest neighbour interpolation
    print("Pixelate")
    # Make a copy which will be saved later on
    pixelated_img = img.copy()
    original_width = pixelated_img.size[0]
    original_height = pixelated_img.size[1]

    pixelated_img = pixelated_img.resize(
        (int(original_width * PIXELATE_SCALING), int(original_height * PIXELATE_SCALING)))

    # scaling back
    pixelated_img = pixelated_img.resize((original_width, original_height), Image.NEAREST)

    save_image('_pixelated', filename, pixelated_img)


# Mosaic helper functions start here
def get_average_RGB(image):
    im = np.array(image)
    w, h, d = im.shape
    return tuple(np.average(im.reshape(w * h, d), axis=0))


def split_image(image, size):
    W, H = image.size[0], image.size[1]
    m = n = size
    w, h = int(W / n), int(H / m)
    imgs = []
    for j in range(m):
        for i in range(n):
            imgs.append(image.crop((i * w, j * h, (i + 1) * w, (j + 1) * h)))
    return imgs


def get_best_match_index(input_avg, avgs):
    avg = input_avg
    index = 0
    min_index = 0
    min_dist = float("inf")
    for val in avgs:
        dist = ((val[0] - avg[0]) * (val[0] - avg[0]) +
                (val[1] - avg[1]) * (val[1] - avg[1]) +
                (val[2] - avg[2]) * (val[2] - avg[2]))
        if dist < min_dist:
            min_dist = dist
            min_index = index
        index += 1
    return min_index


def create_image_grid(images, grid_size):
    m = n = grid_size
    width = max([img.size[0] for img in images])
    height = max([img.size[1] for img in images])
    grid_img = Image.new('RGB', (n * width, m * height))
    for index in range(len(images)):
        row = int(index / n)
        col = index - n * row
        grid_img.paste(images[index], (col * width, row * height))
    return grid_img


def create_photomosaic(target_image, input_images, grid_size,
                       reuse_images=True):
    target_images = split_image(target_image, grid_size)

    output_images = []
    count = 0
    # For progress printer only
    # batch_size = int(len(target_images) / 10)
    avgs = []
    for img in input_images:
        try:
            avgs.append(get_average_RGB(img))
        except ValueError:
            continue

    for img in target_images:
        avg = get_average_RGB(img)
        match_index = get_best_match_index(avg, avgs)
        output_images.append(input_images[match_index])
        # Progress printer
        # if count > 0 and batch_size > 10 and count % batch_size == 0:
        # print('processed %d of %d...' % (count, len(target_images)))
        count += 1
        # remove selected image from input if flag set
        if not reuse_images:
            input_images.remove(match_index)

    mosaic_image = create_image_grid(output_images, grid_size)
    return mosaic_image


def get_images(images_directory):
    files = os.listdir(images_directory)
    images = []

    # Having a maximum number of distinct images creating the mosaic
    # otherwise 1000 images has to be loaded
    i = 0
    for file in files:
        file_path = os.path.abspath(os.path.join(images_directory, file))
        try:
            fp = open(file_path, "rb")
            im = Image.open(fp)
            images.append(im)
            im.load()
            fp.close()
        except Exception as base_ex:
            print(base_ex)
        i += 1

        if i >= 10:
            break
    return images


# Mosaic helper functions end here


def mosaic(base_path, img, filename):
    # Having it global just not to load many times, only once / run
    global INPUT_IMG_LIST

    if not INPUT_IMG_LIST:
        INPUT_IMG_LIST = get_images(base_path)
    # Make a copy which will be saved later on
    mosaic_img = img.copy()
    print("Mosaic")
    if RESIZE_INPUT:
        # print('resizing image...')
        # for given grid size, compute max dims w,h of tiles
        dims = (int(mosaic_img.size[0] / GRID_SIZE),
                int(mosaic_img.size[1] / GRID_SIZE))
        # print("max tile dims: %s" % (dims,))
        # resize
        for img in INPUT_IMG_LIST:
            img.thumbnail(dims)
    # Make a copy which will be saved later on
    mosaic_img = create_photomosaic(mosaic_img, INPUT_IMG_LIST, GRID_SIZE, REUSE_IMAGES)

    save_image('_mosaic', filename, mosaic_img)


def transform_images(base_path, img_name):
    # Opens a image in RGB mode
    with Image.open(base_path + img_name) as image:
        for transformer in transformation_functions:
            try:
                if 'mosaic' == transformer.__qualname__:
                    transformer(base_path, image, img_name)
                else:
                    transformer(image, img_name)
            except Exception as base_ex:
                print(base_ex)


def transform_files():
    # load files one-by-one
    print("[INFO] loading images from source: {}".format(ORIGINAL_FILE_PATH))
    i = 0
    for filename in os.listdir(ORIGINAL_FILE_PATH):

        if filename.endswith(".png") or filename.endswith(".jpg") or filename.endswith(".jpeg"):
            i += 1
            print("[INFO] transforming {}.image: {}".format(i, filename))
            transform_images(ORIGINAL_FILE_PATH, filename)
            continue
        else:
            print("TBD")
            continue


transformation_functions = [image_enhance, mosaic, pixelate, invert, crop, flip, stretch,
                            base_filter, edges, insta_filter]

if __name__ == "__main__":
    # Collect input / output folders with having the
    # option to have command line arguments or constants
    parse_folders()
    transform_files()
