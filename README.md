# Eco-Wise - Reverse Vending Machine for Plastic Bottle Recycling

## Overview

Eco-Wise is an innovative reverse vending machine system designed to encourage recycling in urban cities. The system allows users to recycle plastic bottles in exchange for rewards. By combining computer vision with a Go-based backend and React frontend, Eco-Wise provides a seamless, user-friendly experience for users while effectively managing and processing the recycling operations.

This project utilizes advanced technologies such as Object Detection with YOLO (You Only Look Once) to recognize plastic bottles and a backend system to track the number of bottles and handle user interactions.

## Features

* **Plastic Bottle Detection:** Uses computer vision with YOLO to detect and count plastic bottles being recycled.
* **Reward System:** For every bottle recycled, users earn rewards, incentivizing recycling and sustainability.
* **Real-time Feedback:** Users get instant feedback on the amount of money or reward points earned from their recycled bottles.
* **User-Friendly Interface:** Simple and intuitive web interface built with React.js for easy interaction.
* **Backend Integration:** A robust Go-based backend that handles real-time data processing, payment handling, and communication with the Arduino hardware via serial communication.

## Project Structure
```
├── .env
├── .gitignore
├── go.mod
├── go.sum
├── .vscode/
│   ├── c_cpp_properties.json
│   ├── launch.json
│   └── settings.json
├── cmd/
│   └── main.go                 # Go application entry point
├── computer_vision/
│   ├── arduino.cpp             # Arduino hardware control
│   ├── camera.py              # Bottle detection script
│   ├── yolov3.weights          # YOLO model weights
│   ├── yolov3.cfg              # YOLO config file
│   └── ... (other CV assets)
├── frontend/
│   ├── public/                 # Static assets
│   ├── src/                    # React components
│   ├── package.json            # Frontend dependencies
│   └── ... (React configuration files)
├── internal/
│   └── serial/
│       └── serial.go           # Arduino serial communication
```
## Setup & Installation

### Prerequisites

* **Go:** The Go backend application is built using Go. Ensure you have Go installed (version 1.16 or above).
* **Node.js and npm:** Required for the frontend React app.
* **Python:** The computer vision system uses Python with libraries such as OpenCV.
* **Arduino:** The system uses an Arduino board to manage bottle detection hardware.

### Backend Setup

1.  **Clone the repository:**
    ```bash
    git clone [https://github.com/abenezer54/eco_wise.git](https://github.com/abenezer54/eco_wise.git)
    cd eco_wise
    ```
2.  **Install Go dependencies:**
    ```bash
    go mod tidy
    ```
3.  **Set up environment variables** by creating a `.env` file in the root directory and adding the following:
    ```plaintext
    TWILIO_ACCOUNT_SID=your_twilio_account_sid
    TWILIO_AUTH_TOKEN=your_twilio_auth_token
    ```
4.  **Run the Go server:**
    ```bash
    go run cmd/main.go
    ```

### Frontend Setup

1.  **Navigate to the frontend directory:**
    ```bash
    cd frontend
    ```
2.  **Install frontend dependencies:**
    ```bash
    npm install
    ```
3.  **Run the frontend development server:**
    ```bash
    npm run dev
    ```
    This will launch the React app at http://localhost:3000.

### Computer Vision Setup

The computer vision system uses YOLO for detecting plastic bottles.

1.  **Install Python dependencies:**
    ```bash
    pip install -r requirements.txt
    ```
2.  **Ensure you have the correct YOLO model files** (`frozen_inference_graph.pb`, `yolov3.cfg`, and `yolov3.weights`) in the `computer_vision/` directory.
3.  **Run the Python camera script** to start detecting bottles:
    ```bash
    python computer_vision/camera.py
    ```

### Serial Communication

Ensure that your Arduino is connected and the correct serial port is specified in the `cmd/main.go` file. You may need to modify the port name (e.g., `COM9`).

## How It Works

* **Bottle Detection:** The system uses a camera and Python script (`camera.py`) to detect plastic bottles via computer vision. The YOLO model processes video frames to identify bottles.
* **Backend:** The Go backend handles the business logic, including the bottle count and payment processing. It communicates with the Arduino to control the physical machine.
* **Frontend:** The React app provides a user interface for users to interact with the system. Users can view their recycled bottles, track rewards, and receive confirmation messages.
* **Rewards System:** The backend manages the bottle count and rewards. When a user recycles, they earn reward points or money, which are processed and sent via Twilio SMS for notifications.

## Usage

* **Recycling:** Place a plastic bottle in front of the machine. The camera detects it, and the backend increments the bottle count.
* **Reward:** After a certain number of bottles, users can redeem their points for rewards (money or other incentives).

## Technologies Used

* **Go (Golang):** For the backend and real-time data processing.
* **React.js:** For the user interface.
* **Python:** For the computer vision bottle detection using YOLO.
* **Twilio:** For sending SMS notifications to users when they receive their rewards.
* **Arduino:** To control the hardware interface of the reverse vending machine.

## License

This project is licensed under the MIT License.
