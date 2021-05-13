import torch
from torchvision import models
from torch.hub import load_state_dict_from_url

# We will apply t-SNE to the features extracted by ResNet101 network. First, let’s discuss how Neural Nets process the data.

# We need to re-implement ResNet to be able to extract the last feature map before the classifier head. To do that, we’ll
# create a class that inherits the standard torchvision ResNet and runs the inference without the final classifier.
# It’s important that this way we are able to load the ImageNet pretrained weights to the network
#Basically: getting it as pre-trained weights

# The usual way of training a network:
#
# You want to train a neural network to perform a task (e.g. classification)
# on a data set (e.g. a set of images). You start training by initializing the
# weights randomly. As soon as you start training, the weights are changed in order
# to perform the task with less mistakes (i.e. optimization). Once you're satisfied ' \
#         'with the training results you save the weights of your network somewhere.
#
# You are now interested in training a network to perform a new task (e.g. object detection)
# on a different data set (e.g. images too but not the same as the ones you used before).
# Instead of repeating what you did for the first network and start from training with randomly
# initialized weights, you can use the weights you saved from the previous network as the initial
# weight values for your new experiment. Initializing the weights this way is referred to as using
# a pre-trained network. The first network is your pre-trained network. The second one is the network you are fine-tuning.
#
# The idea behind pre-training is that random initialization is...well...random, the values of the weights have nothing
# to do with the task you're trying to solve. Why should a set of values be any better than another set? But how else would ' \
#         'you initialize the weights? If you knew how to initialize them properly for the task, you might as well set them ' \
#         'to the optimal values (slightly exaggerated). No need to train anything. You have the optimal solution to your problem.' \
#         ' Pre-training gives the network a head start. As if it has seen the data before.

# Define the architecture by modifying resnet.
# Original code is here
# https://github.com/pytorch/vision/blob/b2e95657cd5f389e3973212ba7ddbdcc751a7878/torchvision/models/resnet.py
class ResNet101(models.ResNet):
    def __init__(self, num_classes=1000, pretrained=True, **kwargs):
        # Start with standard resnet101 defined here
        # https://github.com/pytorch/vision/blob/b2e95657cd5f389e3973212ba7ddbdcc751a7878/torchvision/models/resnet.py
        super().__init__(block=models.resnet.Bottleneck, layers=[3, 4, 23, 3], num_classes=num_classes, **kwargs)
        if pretrained:
            state_dict = load_state_dict_from_url(models.resnet.model_urls['resnet101'], progress=True)
            self.load_state_dict(state_dict)

    # Reimplementing forward pass.
    # Replacing the following code
    # https://github.com/pytorch/vision/blob/b2e95657cd5f389e3973212ba7ddbdcc751a7878/torchvision/models/resnet.py#L197-L213
    def _forward_impl(self, x):
        # Standard forward for resnet
        x = self.conv1(x)
        x = self.bn1(x)
        x = self.relu(x)
        x = self.maxpool(x)

        x = self.layer1(x)
        x = self.layer2(x)
        x = self.layer3(x)
        x = self.layer4(x)

        # Notice there is no forward pass through the original classifier.
        x = self.avgpool(x)
        x = torch.flatten(x, 1)

        return x