#include "lcd-display.hh"

LcdDisplay::LcdDisplay() { }

LcdDisplay::LcdDisplay(const int displayWidth, const int displayHeight, const int pinRegisterSelect, const int pinEnable, const int pinDataBus1, const int pinDataBus2, const int pinDataBus3, const int pinDataBus4) {
    if (displayWidth <= 0) {
        throw std::runtime_error("LcdDisplay: Invalid display width specified.");
    }

    this.mDisplayWidth = displayWidth;

    if (displayHeight <= 0) {
        throw std::runtime_error("LcdDisplay: Invalid display height specified.");
    }

    this.mDisplayHeight = displayHeight;

    if (pinRegisterSelect <= 0) {
        throw std::runtime_error("LcdDisplay: Invalid register selection pin specified.");
    }

    if (pinEnable <= 0) {
        throw std::runtime_error("LcdDisplay: Invalid display enable pin specified.");
    }

    if (pinDataBus1 <= 0) {
        throw std::runtime_error("LcdDisplay: Invalid data bus 1 pin specified.");
    }

    if (pinDataBus2 <= 0) {
        throw std::runtime_error("LcdDisplay: Invalid data bus 2 pin specified.");
    }

    if (pinDataBus3 <= 0) {
        throw std::runtime_error("LcdDisplay: Invalid data bus 3 pin specified.");
    }

    if (pinDataBus4 <= 0) {
        throw std::runtime_error("LcdDisplay: Invalid data bus 4 pin specified.");
    }

    this.mpLcdDisplayDevice = new LiquidCrystal(pinRegisterSelect, pinEnable, pinDataBus1, pinDataBus2, pinDataBus3, pinDataBus4);
    this.mpLcdDisplayDevice->begin(this.mDisplayWidth, this.mDisplayHeight);
    this.mpLcdDisplayDevice->clear();
}

LcdDisplay::~LcdDisplay() {
    delete this.mpLcdDisplayDevice;
}

void LcdDisplay::clear() {
    this.mpLcdDisplayDevice->clear();
}

void LcdDisplay::writeLine(int lineIndex, String textValue) {
    if (lineIndex < 0 || lineIndex >= this.mDisplayHeight) {
        throw std::runtime_error("LcdDisplay: Invalid line index specified to display text.");
    }

    String padding(this.mDisplayWidth - textValue.length(), ' ');

    char* lineBuffer = new char(this.mDisplayWidth);
    snprintf(lineBuffer, this.mDisplayWidth, "%s%s", textValue, padding);

    this.mpLcdDisplayDevice->setCursor(0, lineIndex);
    this.mpLcdDisplayDevice->print(lineBuffer);

    delete[] lineBuffer;
}
