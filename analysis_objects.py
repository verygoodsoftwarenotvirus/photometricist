import base64
from PIL import Image
from requests import get


class Product():
    def __init__(self, sku=None):
        self.sku = sku
        self.title = None
        self.photo_url = None
        self.photo = None

    def retrieve_photo(self, dest_folder):
        target_file_name = "{0}{1}".format(dest_folder,
                                           self.photo_url[self.photo_url.rfind("/"):self.photo_url.rfind("?")])
        data = get(self.photo_url)
        with open(target_file_name, 'wb') as f:
            for chunk in data.iter_content():
                f.write(chunk)
        self.photo = Photo(target_file_name)


class Photo():
    def __init__(self, filename):
        self.image = Image.open(filename)
        self.original_path = filename

        self.cropped_path = None
        self.cropped_width = None

        self.computed_colors = []
        self.strategy_color_matches = {}
        self.base64_encoding = None

    def base64_encode(self, image_type="cropped"):
        if image_type.lower() == "cropped":
            path = self.cropped_path
        elif image_type.lower() == "original":
            path = self.original_path
        else:
            raise ValueError("Given image type {0} is invalid.".format(image_type))

        photo_extension = path[path.rfind(".") + 1:]
        if photo_extension == "jpg":
            photo_extension = "jpeg"
        b64_prefix = "data:image/{0};base64,".format(photo_extension)

        with open(path, "rb") as image:
            # encoding adds some junk on the sides.
            result = str(base64.b64encode(image.read()))[2:-1]

        self.base64_encoding = "{0}{1}".format(b64_prefix, result)

    def create_photo_thumbnail(self, image, size=300):
        self.image = image.thumbnail((size, size))

    def crop_and_save_photo(self, crop_percentage, photo_folder, cropped_folder):
        self.image = self.center_crop_image_by_percentage(crop_percentage)
        self.cropped_path = self.original_path.replace(photo_folder, cropped_folder)
        self.image.save(self.cropped_path)
        self.cropped_width = Image.open(self.cropped_path).size[0]

    def center_crop_image_by_percentage(self, percentage=0):
        percentage = int(min(percentage, 100))
        modifier = 0
        if percentage:
            modifier = percentage * .01

        width, height = self.image.size
        left = abs(int(width - (width * modifier)))
        right = abs(int(height - (height * modifier)))
        top = abs(int(width * modifier))
        bottom = abs(int(height * modifier))

        return self.image.crop((left, right, top, bottom))
