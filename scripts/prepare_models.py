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

model1 = applications.EfficientNetB7(weights='imagenet', include_top=False, pooling='avg')
model2 = applications.EfficientNetB6(weights='imagenet', include_top=False, pooling='avg')
model3 = applications.inception_resnet_v2.InceptionResNetV2(weights='imagenet', include_top=False, pooling='avg')
model4 = applications.DenseNet201(weights='imagenet', include_top=False, pooling='avg')
model5 = applications.inception_v3.InceptionV3(weights='imagenet', include_top=False, pooling='avg')
model6 = applications.NASNetLarge(weights='imagenet', include_top=False, pooling='avg')
model7 = applications.ResNet152V2(weights='imagenet', include_top=False, pooling='avg')


img = image.load_img('Ardee_Arollado__Number_08.png', target_size=(224, 224))
x = image.img_to_array(img) # convert image to numpy array
x = np.expand_dims(x, axis=0) # the image is now in an array of shape (3, 224, 224) but we need to expand it to (1, 2, 224, 224) as Keras is expecting a list of images
x = preprocess_input(x)
x = tf.convert_to_tensor(x, dtype=tf.float32)

preds1 = model1.predict(x)[0]
preds2 = model2.predict(x)[0]
preds3 = model3.predict(x)[0]
preds4 = model4.predict(x)[0]
preds5 = model5.predict(x)[0]
preds6 = model6.predict(x)[0]
preds7 = model7.predict(x)[0]

fingerprint_vector = prepare_image_fingerprint_data_for_export_func(preds1)
print('Prediction1:', fingerprint_vector)

fingerprint_vector = prepare_image_fingerprint_data_for_export_func(preds2)
print('Prediction2:', fingerprint_vector)

fingerprint_vector = prepare_image_fingerprint_data_for_export_func(preds3)
print('Prediction3:', fingerprint_vector)

fingerprint_vector = prepare_image_fingerprint_data_for_export_func(preds4)
print('Prediction4:', fingerprint_vector)

fingerprint_vector = prepare_image_fingerprint_data_for_export_func(preds5)
print('Prediction5:', fingerprint_vector)

fingerprint_vector = prepare_image_fingerprint_data_for_export_func(preds6)
print('Prediction6:', fingerprint_vector)

fingerprint_vector = prepare_image_fingerprint_data_for_export_func(preds7)
print('Prediction7:', fingerprint_vector)

model1.save('EfficientNetB7.tf')
model2.save('EfficientNetB6.tf')
model3.save('InceptionResNetV2.tf')
model4.save('DenseNet201.tf')
model5.save('InceptionV3.tf')
model6.save('NASNetLarge.tf')
model7.save('ResNet152V2.tf')
