from flask import Flask, request, jsonify
from PIL import Image, ImageOps, ImageFilter
import pytesseract
import io

app = Flask(__name__)

@app.post("/ocr")
def ocr():
    if "image" not in request.files:
        return jsonify({"error": "missing image"}), 400
    company = (request.form.get("company") or "").strip().lower()
    image_file = request.files["image"]
    try:
        image = Image.open(io.BytesIO(image_file.read()))
        image = ImageOps.exif_transpose(image)
        image = ImageOps.autocontrast(image)
        image = image.filter(ImageFilter.UnsharpMask(radius=1, percent=150, threshold=3))

        text_full = pytesseract.image_to_string(image)
        text_specs = ""
        text_sow = ""
        text_name = ""

        variety_text = ""
        emerge_text = ""
        depth_text = ""
        seed_spacing_text = ""
        row_spacing_text = ""
        thinning_text = ""
        dtm_text = ""

        if company == "botanical_interests":
            try:
                osd = pytesseract.image_to_osd(image)
                rotate = 0
                for line in osd.splitlines():
                    if "Rotate:" in line:
                        rotate = int(line.split(":")[1].strip())
                        break
                if rotate:
                    image = image.rotate(360 - rotate, expand=True)
            except Exception:
                pass

            width, height = image.size

            boxes = {
                "name": {"x": 0.099, "y": 0.091, "w": 0.25, "h": 0.103},
                "variety": {"x": 0.099, "y": 0.188, "w": 0.267, "h": 0.099},
                "emerge": {"x": 0.132, "y": 0.558, "w": 0.235, "h": 0.036},
                "depth": {"x": 0.135, "y": 0.611, "w": 0.235, "h": 0.039},
                "seed_spacing": {"x": 0.135, "y": 0.674, "w": 0.239, "h": 0.06},
                "row_spacing": {"x": 0.137, "y": 0.753, "w": 0.241, "h": 0.037},
                "thinning": {"x": 0.135, "y": 0.812, "w": 0.242, "h": 0.062},
                "days_to_maturity": {"x": 0.137, "y": 0.893, "w": 0.246, "h": 0.039},
                "sow": {"x": 0.375, "y": 0.505, "w": 0.598, "h": 0.283},
            }

            def crop_box(box):
                left = int(width * box["x"])
                top = int(height * box["y"])
                right = int(width * (box["x"] + box["w"]))
                bottom = int(height * (box["y"] + box["h"]))
                return image.crop((left, top, right, bottom))

            def ocr_box(box, psm="6"):
                return pytesseract.image_to_string(crop_box(box), config=f"--oem 1 --psm {psm}")

            text_name = ocr_box(boxes["name"]).strip()
            text_sow = ocr_box(boxes["sow"]).strip()

            specs_lines = []
            variety_text = ocr_box(boxes["variety"]).strip()
            if variety_text:
                specs_lines.append(f"Variety: {variety_text}")
            emerge_text = ocr_box(boxes["emerge"]).strip()
            if emerge_text:
                specs_lines.append(f"Days to emerge: {emerge_text}")
            depth_text = ocr_box(boxes["depth"]).strip()
            if depth_text:
                specs_lines.append(f"Sow depth: {depth_text}")
            seed_spacing_text = ocr_box(boxes["seed_spacing"]).strip()
            if seed_spacing_text:
                specs_lines.append(f"Seed spacing: {seed_spacing_text}")
            row_spacing_text = ocr_box(boxes["row_spacing"]).strip()
            if row_spacing_text:
                specs_lines.append(f"Row spacing: {row_spacing_text}")
            thinning_text = ocr_box(boxes["thinning"]).strip()
            if thinning_text:
                specs_lines.append(f"Thinning: {thinning_text}")
            dtm_text = ocr_box(boxes["days_to_maturity"]).strip()
            if dtm_text:
                specs_lines.append(f"Days to maturity: {dtm_text}")

            text_specs = "\n".join(specs_lines)
    except Exception as exc:
        return jsonify({"error": "ocr failed", "detail": str(exc)}), 500

    def clean_text(value):
        return " ".join(value.split()).strip()

    return jsonify({
        "text": text_full,
        "text_specs": text_specs,
        "text_sow": text_sow,
        "text_name": text_name,
        "name": clean_text(text_name),
        "variety": clean_text(variety_text),
        "days_to_emerge": clean_text(emerge_text),
        "sow_depth": clean_text(depth_text),
        "seed_spacing": clean_text(seed_spacing_text),
        "row_spacing": clean_text(row_spacing_text),
        "thinning": clean_text(thinning_text),
        "days_to_maturity": clean_text(dtm_text),
        "when_to_sow": clean_text(text_sow)
    })

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=9010)
