#include <Servo.h>
#include <LiquidCrystal.h>

Servo myServo;
const int servoPin = 13;
int servoPos = 90;

const int rs = 7, en = 12, d4 = 8, d5 = 9, d6 = 10, d7 = 11;
LiquidCrystal lcd(rs, en, d4, d5, d6, d7);

const int continueButtonPin = 2;
const int endButtonPin = 3;

int bottleCount = 0;
const int rewardPerBottle = 5;
bool waitingForButton = true;

void setup() {
  myServo.attach(servoPin);
  myServo.write(servoPos);
  Serial.begin(9600);
  lcd.begin(16, 2);
  displayIntro();
  delay(3000);
  displayStartMessage();
  pinMode(continueButtonPin, INPUT_PULLUP); 
  pinMode(endButtonPin, INPUT_PULLUP);
}

void loop() {
  if (!waitingForButton) {
    if (Serial.available() > 0) {
      String command = Serial.readStringUntil('\n');
      command.trim();


      if (command == "accept") {
        lcd.clear();
        lcd.print("Processing...");
        moveLeft();
        Serial.print("accept");
        lcd.clear();
        lcd.print("Thank you!");
        lcd.setCursor(0, 1);
        lcd.print("Bottle Count: ");
        lcd.print(bottleCount + 1);
        delay(1500);
        bottleCount++;
        waitingForButton = true;
        displayWaitForButtonMessage();
      }
      else if (command == "reject"){
        lcd.clear();
        lcd.print("Processing...");
        moveRight();
        lcd.clear();
        lcd.print("Thank you!");
        lcd.setCursor(0, 1);
        lcd.print("Bottle Count: ");
        lcd.print(bottleCount);
        delay(1500);
        waitingForButton = true;
        displayWaitForButtonMessage();
      }
    }
  } else {
    if (digitalRead(continueButtonPin) == LOW) {
      delay(200); // Debounce
      waitingForButton = false;
      displayStartMessage();
    } else if (digitalRead(endButtonPin) == LOW) {
      // delay(200); // Debounce
      // endProcess();
      // serial.print("finish");
      // waitingForButton = true;
      // displayStartMessage();
    }
  }
}

// ...existing code...

void moveRight() {
  // Start at 0
//   myServo.write(0);
//   servoPos = 0;
  delay(300);

  // Move up to 180
  for (servoPos = 90; servoPos <= 180; servoPos++) {
    myServo.write(servoPos);
    delay(15);
  }

  delay(1000);

  // Move back down to 0
  for (servoPos = 180; servoPos > 90; servoPos--) {
    myServo.write(servoPos);
    delay(15);
  }

  // End at 0
//   myServo.write(0);
//   servoPos = 0;
}

void moveLeft() {
  // Start at 0
//   myServo.write(0);
//   servoPos = 0;
  delay(300);

  // Move up to 180
  for (servoPos = 90; servoPos > 0; servoPos--) {
    myServo.write(servoPos);
    delay(15);
  }

  delay(1000);

  // Move back down to 0
  for (servoPos = 0; servoPos < 90; servoPos++) {
    myServo.write(servoPos);
    delay(15);
  }

  delay(1000);


}




void endProcess() {
  lcd.clear();
  lcd.print("Recycling Done!");
  lcd.setCursor(0, 1);
  lcd.print("Total: ");
  lcd.print(bottleCount);
  delay(2000);
  lcd.clear();
  lcd.print("Reward Earned:");
  lcd.setCursor(0, 1);
  lcd.print(bottleCount * rewardPerBottle);
  lcd.print(" Units");
  delay(3000);
  bottleCount = 0; // Reset for next session
}



void displayIntro() {
  lcd.clear();
  lcd.print("Welcome to");
  lcd.setCursor(0, 1);
  lcd.print("Smart Recycling");
}

void displayStartMessage() {
  lcd.clear();
  lcd.print("Ready for");
  lcd.setCursor(0, 1);
  lcd.print("Next Bottle!");
}

void displayWaitForButtonMessage() {
  lcd.clear();
  lcd.print("Press Continue");
  lcd.setCursor(0, 1);
  lcd.print("or End");
}

