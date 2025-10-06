from PIL import Image, ImageEnhance, ImageDraw
import matplotlib.pyplot as plt

# Load the original image
image_path = "your_image.png"  # Replace with your image file name
original_image = Image.open(image_path).convert("RGBA")

# Enhance contrast and brightness
enhancer = ImageEnhance.Contrast(original_image)
image_contrasted = enhancer.enhance(1.3)
enhancer = ImageEnhance.Brightness(image_contrasted)
image_brightened = enhancer.enhance(1.2)

# Create a new layer to draw colors
colored_layer = Image.new("RGBA", original_image.size)
draw = ImageDraw.Draw(colored_layer)

# Add some fun colors to different zones
draw.ellipse((120, 50, 280, 210), fill=(0, 0, 0, 120))  # Hair
draw.rectangle((140, 180, 240, 220), fill=(255, 0, 0, 120))  # Glasses
draw.rectangle((170, 250, 260, 400), fill=(220, 220, 220, 100))  # Skull base
draw.ellipse((40, 30, 100, 90), fill=(255, 105, 180, 100))  # Left rose
draw.ellipse((300, 30, 360, 90), fill=(186, 85, 211, 100))  # Right rose
draw.rectangle((260, 320, 300, 400), fill=(139, 69, 19, 120))  # Guitar body

# Combine the image and the color layer
colored_image = Image.alpha_composite(image_brightened, colored_layer)

# Show the final result
plt.figure(figsize=(8, 10))
plt.imshow(colored_image)
plt.axis("off")
plt.title("Digitally Colored Drawing")
plt.show()
