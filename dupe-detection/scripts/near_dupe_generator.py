# import the necessary packages
import argparse
import json
import os
import random
from random import randrange
from PIL import Image, ImageOps, ImageFilter, ImageEnhance
import PIL
import numpy as np
import hashlib

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
CROP_RATIO = random.uniform(0.01, 0.5)
# Stretch width ratio used by random_stretch function
STRETCH_WIDTH_RATIO = random.uniform(0.1, 2.0)
# Stretch height ratio used by random_stretch function
STRETCH_HEIGHT_RATIO = random.uniform(0.1, 2.0)
# Radius parameter for GaussianBlur
BLUR_RADIUS = randrange(1, 30)
# Determines the scaling factor for step1 in pixelation
PIXELATE_SCALING = random.uniform(0.05, 0.2)
# Enhancement ratio: (min 0, max 1)
# Factor 1.0 always returns a copy of the original image.
# Lower factors mean less color (brightness, contrast, etc).
ENHANCEMENT_FACTOR = random.uniform(0.01, 10.0)
# Mosaic: resize the input to fit original image size?
RESIZE_INPUT = True # Not recommended to have this parameter random !!!!
# Mosaic: grid size. THe bigger this number the the more detailed the image (and slower the process)
GRID_SIZE = randrange(50, 150)
# Mosaic: re-use any image in input
REUSE_IMAGES = True
# Mosaic: List of images which will be contained in the mosaic
INPUT_IMG_LIST = []
# Store every kind of information related to transformations
# and store it in JSON format
LOGGER = None
# Maximum nr. of transformations/image
MAXIMUM_NR_OF_TRANSFORMATIONS = 4


class JsonLogger:
    """Logger class for JSON export"""
    root_info = {}
    asset_level = {}

    def __init__(self):
        pass

    def add_element(self, img_name, transformation_info):
        """Add an element to the root dictionary
        Parameters
        ----------
        img_name : str
            The name of the image
        transformation_info : dict (nested dictionary)
            transformation_info: containing: hashes (input, output), transformation parameters, and informations
        """
        self.asset_level[img_name] = transformation_info

    def add_top_lvl(self):
        self.root_info['assets'] = self.asset_level


# Function to create random number of iterations, transformations
def randomizer(lower_incl, upper_excl):
    irand = randrange(lower_incl, upper_excl)

    return irand


def randomize_global_params():
    # Crop size ratio used by random_crop function
    global CROP_RATIO
    CROP_RATIO = random.uniform(0.01, 0.5)
    # Stretch width ratio used by random_stretch function
    global STRETCH_WIDTH_RATIO
    STRETCH_WIDTH_RATIO = random.uniform(0.1, 2.0)
    # Stretch height ratio used by random_stretch function
    global STRETCH_HEIGHT_RATIO
    STRETCH_HEIGHT_RATIO = random.uniform(0.1, 2.0)
    # Radius parameter for GaussianBlur
    global BLUR_RADIUS
    BLUR_RADIUS = randrange(1, 30)
    # Determines the scaling factor for step1 in pixelation
    global PIXELATE_SCALING
    PIXELATE_SCALING = random.uniform(0.05, 0.2)
    # Enhancement ratio: (min 0, max 1)
    # Factor 1.0 always returns a copy of the original image.
    # Lower factors mean less color (brightness, contrast, etc).
    global ENHANCEMENT_FACTOR
    ENHANCEMENT_FACTOR = random.uniform(0.01, 10.0)
    # Mosaic: grid size. THe bigger this number the the more detailed the image (and slower the process)
    global GRID_SIZE
    GRID_SIZE = randrange(50, 150)


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


def get_hash(filename) -> object:
    print("Filename: {}".format(filename))

    with open(filename, "rb") as f:
        bytes = f.read()
        calculated_hash = hashlib.sha256(bytes).hexdigest();
        print(calculated_hash)

    return calculated_hash


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
def save_image(transform_type, filename, to_be_saved, img_transformation_dict):
    file = os.path.splitext(filename)[0]
    file_extension = os.path.splitext(filename)[1]

    new_filename = file + transform_type + file_extension
    to_be_saved.save(TRANSFORMED_FILE_PATH + new_filename)
    to_be_saved.close()
    img_transformation_dict['transformed_img_name'] = new_filename
    return get_hash(TRANSFORMED_FILE_PATH + new_filename)


def crop(cropped_img, filename):
    print("Cropping image")
    param_dictionary = {}
    # get image size
    x, y = cropped_img.size

    # determine the mask for cropping
    # now i pick up the 30% of the width
    # and 30% of the height
    matrix_x = int(x * CROP_RATIO)
    matrix_y = int(y * CROP_RATIO)
    x1 = randrange(0, x - matrix_x)
    y1 = randrange(0, y - matrix_y)

    param_dictionary['crop_ratio'] = CROP_RATIO
    param_dictionary['crop_left_coord'] = x1
    param_dictionary['crop_upper_coord'] = y1
    param_dictionary['crop_right_coord'] = (x1 + matrix_x)
    param_dictionary['crop_lower_coord'] = (y1 + matrix_y)
    try:
        cropped_img = cropped_img.crop((x1, y1, x1 + matrix_x, y1 + matrix_y))
    except Exception as bs:
        print("Exception caught during cropping! Issue: {}".format(bs))
        return None, None

    # No saving of individual image required. Left it here maybe future need..
    # save_image('_cropped', filename, cropped_img)
    return cropped_img, param_dictionary


def flip(flipped_img, filename):
    print("Flipping image")
    param_dictionary = {}
    # random nr. between horizontally vs. vertically flipping
    irand = randrange(0, 2)
    if irand == 0:
        # flip image vertically
        param_dictionary['flipping_direction'] = 'Vertical'
    else:
        # flip image horizontally
        param_dictionary['flipping_direction'] = 'Horizontal'

    try:
        flipped_img = flipped_img.transpose(irand)
    except Exception as bs:
        print("Exception caught during flipping! Issue: {}".format(bs))
        return None, None

    # No saving of individual image required. Left it here maybe future need..
    # save_image('_flipped', filename, flipped_img)
    return flipped_img, param_dictionary


def stretch(stretched_img, filename):
    print("Stretching image")
    param_dictionary = {}
    # get image size
    x, y = stretched_img.size

    new_width = int(x * STRETCH_WIDTH_RATIO)
    new_height = int(y * STRETCH_HEIGHT_RATIO)
    try:
        stretched_img = stretched_img.resize((new_width, new_height), Image.ANTIALIAS)
        param_dictionary['stretch_width_ratio'] = STRETCH_WIDTH_RATIO
        param_dictionary['stretch_height_ratio'] = STRETCH_HEIGHT_RATIO
    except Exception as bs:
        print("Exception caught during stretching! Issue: {}".format(bs))
        return None, None

    # No saving of individual image required. Left it here maybe future need..
    # save_image('_stretched', filename, stretched_img)
    return stretched_img, param_dictionary


def base_filter(filtered_img, filename):
    irand = randrange(0, 9)

    param_dict = {}

    filter_type = base_filter_types(irand)
    filter_attribute: object = getattr(ImageFilter, filter_type)
    print("Image is filtered based on: {} base filter".format(filter_type))
    try:
        if filter_type == 'BLUR':
            filtered_img = filtered_img.filter(ImageFilter.GaussianBlur(BLUR_RADIUS))
            param_dict['base_filter_type'] = 'Pillow.GaussianBlur with radius: {}'.format(BLUR_RADIUS)
        else:
            filtered_img = filtered_img.filter(filter_attribute)
            param_dict['base_filter_type'] = 'Pillow.{}'.format(filter_type)
    except Exception as bs:
        print("Exception caught during base filter! Issue: {}".format(bs))
        return None, None

    # No saving of individual image required. Left it here maybe future need..
    # save_image('_base_filtered_{}'.format(filter_type), filename, filtered_img)
    return filtered_img, param_dict


def insta_filter(filtered_img, filename):
    irand = randrange(0, 26)
    param_dict = {}
    filter_type = insta_filter_types(irand)

    # Trying to create/access function from dictionary string
    try:
        func = getattr(pilgram, filter_type)
        filtered_img = func(filtered_img)
        print("Image is filtered based on: {} insta filter".format(filter_type))
        param_dict['insta_filter_type'] = 'Pilgram.{}'.format(filter_type)
    except AttributeError:
        print("Function not found")
        return None, None
    except Exception as bs:
        print("Exception caught during insta filter! Issue: {}".format(bs))
        return None, None

    # No saving of individual image required. Left it here maybe future need..
    # save_image('_insta_filtered_{}'.format(filter_type), filename, filtered_img)

    return filtered_img, param_dict


def edges(filtered_img, filename):
    print("Finding edges")

    param_dict = {'find_edges_algorithm': 'PIL.ImageFilter.FIND_EDGES'}

    try:
        filtered_img = filtered_img.filter(ImageFilter.FIND_EDGES)
    except Exception as bs:
        print("Exception caught during edge function filter! Issue: {}".format(bs))
        return None, None

    # No saving of individual image required. Left it here maybe future need..
    # save_image('_edges', filename, filtered_img)

    return filtered_img, param_dict


def invert(inverted_img, filename):
    print("Invert")

    param_dict = {'invert_algorithm': 'PIL.ImageOps.Invert(img)'}
    try:
        inverted_img = ImageOps.invert(inverted_img)
    except IOError:
        # Convert to RGB
        inverted_img = ImageOps.invert(inverted_img.convert('RGB'))
    except Exception as bs:
        print("Problem in inversion is: {}".format(bs))
        return None, None



    # No saving of individual image required. Left it here maybe future need..
    # save_image('_inverted', filename, inverted_img)

    return inverted_img, param_dict


def image_enhance(saturated_img, filename):
    print("Saturate")
    param_dict = {}
    irand = randrange(0, 3)

    filter_type = enhancement_types(irand)
    filter_attribute = getattr(PIL.ImageEnhance, filter_type)
    print("Image is filtered based on: {} enhance types".format(filter_type))

    try:
        converter = filter_attribute(saturated_img)
        saturated_img = converter.enhance(ENHANCEMENT_FACTOR)
        param_dict['enhancement_type'] = filter_type
        param_dict['enhancement_factor'] = ENHANCEMENT_FACTOR
    except Exception as bs:
        print("Problem during enhancement is: {}".format(bs))
        return None, None

    # No saving of individual image required. Left it here maybe future need..
    # save_image('_enhanced_{}'.format(filter_type), filename, saturated_img)

    return saturated_img, param_dict


def pixelate(pixelated_img, filename):
    # Simply put pixelation is:
    # 1. Resize to a smaller size with the help of PIXELATE_SCALING (smaller the number the more pixelated the result
    # 2. Scale it back to original using nearest neighbour interpolation
    print("Pixelate")
    param_dict = {}

    original_width = pixelated_img.size[0]
    original_height = pixelated_img.size[1]

    param_dict['pixelate_resize_ratio'] = PIXELATE_SCALING
    param_dict['pixelate_resize_to_original_algorithm'] = 'PIL.Image.NEAREST'

    try:

        pixelated_img = pixelated_img.resize(
            (int(original_width * PIXELATE_SCALING), int(original_height * PIXELATE_SCALING)))

        # scaling back
        pixelated_img = pixelated_img.resize((original_width, original_height), Image.NEAREST)
    except Exception as bs:
        print("Problem during pixelation is: {}".format(bs))
        return None, None

    # No saving of individual image required. Left it here maybe future need..
    # save_image('_pixelated', filename, pixelated_img)

    return pixelated_img, param_dict


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


def get_images(images_directory, params_dict):
    files = os.listdir(images_directory)
    images = []
    # Randomly get the maximum of 10 images which will take part in the mosaic creation
    random.shuffle(files)

    # Having a maximum number of distinct images creating the mosaic
    # otherwise 1000 images has to be loaded
    i = 0
    for file in files:
        file_key = 'mosaic_pic_{}'.format(i)
        params_dict[file_key] = file
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


def mosaic(base_path, mosaic_img, filename):
    # Having it global just not to load many times, only once / run
    param_dictionary = {}
    global INPUT_IMG_LIST

    INPUT_IMG_LIST = []

    if not INPUT_IMG_LIST:
        INPUT_IMG_LIST = get_images(base_path, param_dictionary)

    param_dictionary['resize_option'] = RESIZE_INPUT
    print("Mosaic")
    if RESIZE_INPUT:
        # print('resizing image...')
        # for given grid size, compute max dims w,h of tiles
        try:
            dims = (int(mosaic_img.size[0] / GRID_SIZE),
                    int(mosaic_img.size[1] / GRID_SIZE))
            param_dictionary['dimension'] = dims
            # print("max tile dims: %s" % (dims,))
            # resize
            for img in INPUT_IMG_LIST:
                img.thumbnail(dims)
        except Exception as bs:
            print("Problem during mosaic dimension is: {}".format(bs))
            return None, None

    # Make a copy which will be saved later on
    try:
        mosaic_img = create_photomosaic(mosaic_img, INPUT_IMG_LIST, GRID_SIZE, REUSE_IMAGES)
    except Exception as bs:
        print("Problem during mosaic creation is: {}".format(bs))
        return None, None


    # No saving of individual image required. Left it here maybe future need..
    # save_image('_mosaic', filename, mosaic_img)
    return mosaic_img, param_dictionary


# Function for applying transformations per image
def get_random_transformation_name(lower_incl, upper_excl, random_transformation_functions):
    function_name = None
    while True:
        function_idx = randomizer(lower_incl, upper_excl)
        if transformation_functions[function_idx] not in random_transformation_functions:
            # Return the name of the function
            function_name = transformation_functions[function_idx]
            break

    return function_name


def transform_each_image(transformation_nr, img_transformation_dict,
                         image, base_path, img_name):

    # THis will hold the sequence of transformation function names
    random_transformation_functions = []
    for i in range(transformation_nr):
        func_name = get_random_transformation_name(0, len(transformation_functions), random_transformation_functions)
        random_transformation_functions.append(func_name)

    # We need to randomize global parmeters
    randomize_global_params()
    # Here we already have a randomized number and sequence of transformations as well
    i = 0;
    for transformer in random_transformation_functions:
        try:
            image_copy = image.copy()
            i += 1
            print("{}. transformation is: {}".format(i, transformer.__qualname__))
            if 'mosaic' == transformer.__qualname__:
                image_copy, param_dictionary = transformer(base_path, image_copy, img_name)
            else:
                image_copy, param_dictionary = transformer(image_copy, img_name)

            # PARAMS: TO BE IMPLEMENTED
            if image is not None and param_dictionary is not None:
                # Rewrite into the image since it will
                image = image_copy.copy()
                transformation_params = ("transformation_{}_params".format(i))
                transformation_name = ("transformation_{}_name".format(i))
                img_transformation_dict[transformation_name] = transformer.__qualname__
                img_transformation_dict[transformation_params] = param_dictionary
            else:
                i -= 1

        except Exception as base_ex:
            print(base_ex)

    new_hash = save_image('_transformed', img_name, image, img_transformation_dict)
    img_transformation_dict['transformed_img_hash'] = new_hash


def transform_images(base_path, img_name):
    
    # Image related information dictionary where we keep track of every single transformation on them
    img_transformation_dict = dict()

    # Get SHA256 calculated hash
    img_transformation_dict['original_hash'] = get_hash(base_path+img_name)
    
    # Opens a image in RGB mode
    with Image.open(base_path + img_name) as image:

        # Get a random number, how many transformations shall be performed per asset
        transformation_nr = randomizer(1, MAXIMUM_NR_OF_TRANSFORMATIONS+1)

        # Perform transformation
        transform_each_image(transformation_nr, img_transformation_dict,
                             image, base_path, img_name)

        # Add each image to the 'asset' level of the LOGGER dictionary
        LOGGER.add_element(img_name, img_transformation_dict)


def transform_files():
    # Load files one-by-one and transform them
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
    # Create JsonLogger instance to story every transformation
    LOGGER = JsonLogger()
    # It does the main work :)
    transform_files()
    # Add 'asset' child to the root
    LOGGER.add_top_lvl()
    # Save it to file
    with open("transformation.json", "w") as outfile:
        json.dump(LOGGER.root_info, outfile, indent=4)

