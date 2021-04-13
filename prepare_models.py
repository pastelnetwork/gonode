import keras
from keras.applications.nasnet import NASNetMobile
from keras.preprocessing import image
from keras.applications.xception import preprocess_input, decode_predictions
from tensorflow.keras import applications
import numpy as np
import tensorflow as tf
from keras import backend as K

def prepare_image_fingerprint_data_for_export_func(image_feature_data):
    image_feature_data_arr = np.char.mod('%f', image_feature_data) # convert from Numpy to a list of values
    x_data = np.asarray(image_feature_data_arr).astype('float64') # convert image data to float64 matrix. float64 is need for bh_sne
    image_fingerprint_vector = x_data.reshape((x_data.shape[0], -1))
    return image_fingerprint_vector

tf.compat.v1.set_random_seed(1234)

model = applications.NASNetLarge(weights='imagenet', include_top=False, pooling='avg')

img = image.load_img('Ardee_Arollado__Number_08.png', target_size=(224, 224))
x = image.img_to_array(img) # convert image to numpy array
x = np.expand_dims(x, axis=0) # the image is now in an array of shape (3, 224, 224) but we need to expand it to (1, 2, 224, 224) as Keras is expecting a list of images
x = preprocess_input(x)
x = tf.convert_to_tensor(x, dtype=tf.float32)

preds = model.predict(x)[0]
fingerprint_vector = prepare_image_fingerprint_data_for_export_func(preds)
print('Prediction:', fingerprint_vector)

model.save('NASNetLarge.tf')