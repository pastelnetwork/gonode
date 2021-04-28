ECHO ON

saved_model_cli show --all --dir EfficientNetB7.tf>EfficientNetB7.tf/EfficientNetB7.txt
saved_model_cli show --all --dir EfficientNetB6.tf>EfficientNetB6.tf/EfficientNetB6.txt
saved_model_cli show --all --dir InceptionResNetV2.tf>InceptionResNetV2.tf/InceptionResNetV2.txt
saved_model_cli show --all --dir DenseNet201.tf>DenseNet201.tf/DenseNet201.txt
saved_model_cli show --all --dir InceptionV3.tf>InceptionV3.tf/InceptionV3.txt
saved_model_cli show --all --dir NASNetLarge.tf>NASNetLarge.tf/NASNetLarge.txt
saved_model_cli show --all --dir ResNet152V2.tf>ResNet152V2.tf/ResNet152V2.txt
