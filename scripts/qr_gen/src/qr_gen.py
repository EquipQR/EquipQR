from qrcode import QRCode
from qrcode.constants import ERROR_CORRECT_M
from PIL import Image

def generate_qr_code(data: str, filename: str, box_size: int = 10, border: int = 4) -> None:
    qr = QRCode(
        version=None,
        error_correction=ERROR_CORRECT_M,
        box_size=box_size,
        border=border,
    )
    qr.add_data(data)
    qr.make(fit=True)

    img: Image.Image = qr.make_image(fill_color="black", back_color="white")
    img.save(filename)

if __name__ == "__main__":
    generate_qr_code("ed9ba222-80de-4e16-a0b9-efdb74f0ec44", "equipment_qr.png")
