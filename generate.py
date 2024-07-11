from PIL import Image
import numpy as np

# Define image dimensions and color
width, height = 16384, 16384
color = (255, 255, 255, 255)  # White with full opacity

# Create a white RGBA image
image = Image.new("RGBA", (width, height), color)

# Save the image as webp
image_path = "16384x16384_white_rgba.png"
image.save(image_path, "PNG")

image_path
