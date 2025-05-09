import cv2
import numpy as np
import socket
import time

# Load YOLO model
net = cv2.dnn.readNet('yolov3.weights', 'yolov3.cfg')
with open('coco.names', 'r') as f:
    classes = [line.strip() for line in f.readlines()]

camera_index = 2
cap = cv2.VideoCapture(camera_index, cv2.CAP_DSHOW)

if not cap.isOpened():
    print(f"Error: Unable to access camera index {camera_index}.")
    exit()

def send_signal_to_backend(signal):
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
        s.connect(('localhost', 9000))
        s.sendall(signal.encode())

print("Camera is open. Press 'c' to capture an image, 'q' to quit.")

while True:
    ret, frame = cap.read()
    if not ret:
        print("Failed to grab frame")
        break

    cv2.imshow("Eco wise - Plastic Detection", frame)

    key = cv2.waitKey(1) & 0xFF
    if key == ord('c'):
        print("Capturing image for detection...")
        height, width, channels = frame.shape
        blob = cv2.dnn.blobFromImage(frame, 0.00392, (416, 416), (0, 0, 0), True, crop=False)
        
        net.setInput(blob)
        outs = net.forward(net.getUnconnectedOutLayersNames())

        class_ids = []
        confidences = []
        boxes = []
        for out in outs:
            for detection in out:
                scores = detection[5:]
                class_id = np.argmax(scores)
                confidence = scores[class_id]
                if confidence > 0.5:
                    center_x = int(detection[0] * width)
                    center_y = int(detection[1] * height)
                    w = int(detection[2] * width)
                    h = int(detection[3] * height)
                    x = int(center_x - w / 2)
                    y = int(center_y - h / 2)
                    boxes.append([x, y, w, h])
                    confidences.append(float(confidence))
                    class_ids.append(class_id)

        indexes = cv2.dnn.NMSBoxes(boxes, confidences, 0.5, 0.4)
        signal_sent = False
        if len(indexes) > 0:
            for i in indexes.flatten():
                label = str(classes[class_ids[i]])
                if label == "bottle":
                    print("Bottle detected! Sending 'accept' command to backend.")
                    send_signal_to_backend('accept:1')
                    signal_sent = True
                    break
        if not signal_sent:
            print("No bottle detected. Sending 'reject' command to backend.")
            send_signal_to_backend('reject:0')

    elif key == ord('q'):
        print("Quitting...")
        break

cap.release()
cv2.destroyAllWindows()