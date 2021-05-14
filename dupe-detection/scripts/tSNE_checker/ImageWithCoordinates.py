import re

class ImageWithCoordinates:
    """
        A class used to represent an image either original or transformed
        with coordinates and having the option of calculate the Euclidean
        distance between original and transformed images.

        ...

        Instance attributes
        ----------
        image_name   : str
            name or path to the image, which directly identifies it
        x            : int
            origo x coordinate on the t-SNE coordinate system
        y            : int
            origo y coordinate on the t-SNE coordinate system
        isOriginal   : bool
            if image is original or not ( original = not having "transformed" in image_name)
        parentName   : str
            possible parent name, which shall be retrieved from the list of image_names

        Methods
        -------
        _is_original(image_name)
            Fills the instance parameter if original or not based on the name / path
        _get_parent_name(image_name)
            Replacing the "_xx_transformed" text in order to have the parent name
        """
    def __init__(self, image_name, x, y):
        self.imageName = image_name
        self.x = x
        self.y = y
        self.isOriginal = self._is_original(image_name)
        # Keep in mind this part is only relevant for creating dataset
        # in live we will not ahve any information about parent (original)
        self.parentName = None
        if not self.isOriginal:
            self.parentName = self._get_parent_name(image_name)
        self.std_from_avg = None
        pass

    def __str__(self):
        return "Image_name is {}\nX is:{}, Y is:{}\nisOriginal: {}\nparentName: {}\n".format(self.image_name, self.x, self.y, self.isOriginal, self.parentName)

    def _is_original(self, image_name):
        ret_val = True

        if "transformed" in image_name:
            ret_val = False

        return ret_val

    def _get_parent_name(self, image_name):

        parent_name = re.sub('_[0-9]*_transformed', '', image_name)
        return parent_name



class EuclideanDistanceMap:
    """
        A class used to represent a 2D map (i.e.: like a chess table)
        with N X N where N:= number of images.
        Class contains a list of images in order (!) and a
        a numpy array which will contain the euclidean distances relative
        to image order.

        ...

        Attributes
        ----------
        image_list   : list of ImageWithCoordinates
            a list from ImageWithCoordinates
        eu_distances : 2D np array

        Methods
        -------
        get_eu_distance_bw_2_images(image_A, image_B)
            Claculates the Euclidean distance between two given points (by filename)
        get_mean_eu_distance()
            Claculates the average Euclidean distance within all points
        """
    def __init__(self, images, eu_distances):
        self.image_list = images
        self.eu_distances = eu_distances
        self.average_distance = None
        pass

    def get_eu_distance_bw_2_images(self, image_A, image_B):
        ret_val = None
        idx_x = -1
        idx_y = -1

        counter = 0
        for image in self.image_list:
            if (image.imageName == image_A):
                idx_x = counter
            if (image.imageName == image_B):
                idx_y = counter
            counter += 1

        if (idx_x != -1 and idx_y != -1):
            ret_val = self.eu_distances[idx_x][idx_y]

        return ret_val

    def get_mean_eu_distance(self):
        # Multiply by 2 needed because the array has for example 4 elements in case of
        # 2x2 array. So it divides the things by 2 times more it needs
        # And shall not be called the .mean() function twice
        ret_val = self.average_distance
        if ret_val is None:
            ret_val = (self.eu_distances.mean() * 2)
            self.average_distance = ret_val
        else:
            ret_val = self.average_distance

        return ret_val
