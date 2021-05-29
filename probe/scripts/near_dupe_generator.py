# import the necessary packages
import argparse
import os
from random import randrange
from PIL import Image
from PIL import ImageFilter
import PIL

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


def filter_types(i):
    filter_map = {
        0: 'BLUR',
        1: 'CONTOUR',
        2: 'DETAIL',
        3: 'EDGE_ENHANCE',
        4: 'EDGE_ENHANCE_MORE',
        5: 'EMBOSS',
        6: 'FIND_EDGES',
        7: 'SHARPEN',
        8: 'SMOOTH',
        9: 'SMOOTH_MORE'
    }
    return filter_map.get(i, "Invalid")


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
    save_image('_cropped_', filename, cropped_img)


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

    save_image('_flipped_', filename, flipped_img)


def stretch(img, filename):
    print("Stretching image")
    # Make a copy which will be saved later on
    stretched_img = img.copy()

    # get image size
    x, y = stretched_img.size

    new_width = int(x * STRETCH_WIDTH_RATIO)
    new_height = int(y * STRETCH_HEIGHT_RATIO)

    stretched_img = stretched_img.resize((new_width, new_height), Image.ANTIALIAS)

    save_image('_stretched_', filename, stretched_img)


def base_filter(img, filename):
    # Make a copy which will be saved later on
    filtered_img = img.copy()
    irand = randrange(0, 10)

    filter_type = filter_types(irand)
    filter_attribute: object=getattr(ImageFilter, filter_type)
    print("Image is filtered based on: {} filter".format(filter_type))
    filtered_img = filtered_img.filter(filter_attribute)

    save_image('_base_filtered_', filename, filtered_img)


def transform_images(base_path, img_name):
    # Opens a image in RGB mode
    with Image.open(base_path + img_name) as image:
        for transformer in transformation_functions:
            try:
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
            continue


transformation_functions = [crop, flip, stretch, base_filter]

if __name__ == "__main__":
    # Collect input / output folders with having the
    # option to have command line arguments or constants
    parse_folders()
    transform_files()
