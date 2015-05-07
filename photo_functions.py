import base64
from PIL import Image


def open_image(filename):
    return Image.open(filename)


def save_image(image, filename):
    image.save(filename)


def get_image_width(image):
    return image.size[0]


def get_image_height(image):
    return image.size[1]


def get_image_size(image):
    return image.size


def base64_encode_image(filename):
    with open(filename, "rb") as image:
        return base64.b64encode(image.read())


def create_photo_thumbnail(image, size=300):
    return image.thumbnail((size, size))


def crop_and_save_photo(photo_path, crop_percentage, photo_destination_folder, cropped_folder):
    image = open_image(photo_path)
    image = center_crop_image_by_percentage(image, crop_percentage)
    photo_path = photo_path.replace(photo_destination_folder, cropped_folder)
    save_image(image, photo_path)
    return image


def center_crop_image_by_percentage(image, percentage=0):
    percentage = int(min(percentage, 100))
    modifier = 0
    if percentage:
        modifier = percentage * .01

    width, height = image.size
    left = abs(int(width - (width * modifier)))
    right = abs(int(height - (height * modifier)))
    top = abs(int(width * modifier))
    bottom = abs(int(height * modifier))

    return image.crop((left,
                       right,
                       top,
                       bottom))
