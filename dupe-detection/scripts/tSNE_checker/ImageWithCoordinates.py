import re

class ImageWithCoordinates:
    """
        A class used to represent an image either original or transformed
        with coordinates and having the option of calculate the Euclidean
        distance between original and transformed images.

        ...

        Attributes
        ----------
        image_name : str
            name or path to the image, which directly identifies it
        x           : int
            origo x coordinate on the t-SNE coordinate system
        y           : int
            origo y coordinate on the t-SNE coordinate system
        isOriginal  : bool
            if image is original or not ( original = not having "transformed" in image_name)
        parentName  : str
            possible parent name, which shall be retrieved from the list of image_names

        Methods
        -------
        _is_original(image_name)
            Fills the instance parameter if original or not based on the name / path
        _get_parent_name(image_name)
            Replacing the "_xx_transformed" text in order to have the parent name
        """
    def __init__(self, image_name, x, y):
        self.image_name = image_name
        self.x = x
        self.y = y
        self.isOriginal = self._is_original(image_name)
        self.parentName = None
        if self.isOriginal == False:
            self.parentName = self._get_parent_name(image_name)


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
