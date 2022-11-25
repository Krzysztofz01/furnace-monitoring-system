#include "lcd-display.hh"

LcdDisplay::LcdDisplay() { }

LcdDisplay::LcdDisplay(const int displayWidth, const int displayHeight, const int pinRegisterSelect, const int pinEnable, const int pinDataBus1, const int pinDataBus2, const int pinDataBus3, const int pinDataBus4) {
    if (displayWidth <= 0) {
        throw std::runtime_error("LcdDisplay: Invalid display width specified.");
    }

    mDisplayWidth = displayWidth;

    if (displayHeight <= 0) {
        throw std::runtime_error("LcdDisplay: Invalid display height specified.");
    }

    mDisplayHeight = displayHeight;

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

    mpLcdDisplayDevice = new LiquidCrystal(pinRegisterSelect, pinEnable, pinDataBus1, pinDataBus2, pinDataBus3, pinDataBus4);
    mpLcdDisplayDevice->begin(mDisplayWidth, mDisplayHeight);
    mpLcdDisplayDevice->clear();
}

LcdDisplay::~LcdDisplay() {
    delete mpLcdDisplayDevice;
}

void LcdDisplay::clear() {
    mpLcdDisplayDevice->clear();
}

void LcdDisplay::writeLine(int lineIndex, String textValue) {
    if (lineIndex < 0 || lineIndex >= mDisplayHeight) {
        throw std::runtime_error("LcdDisplay: Invalid line index specified to display text.");
    }

    char* outputBuffer = new char[mDisplayWidth + 1];
    for (int buffIndex = 0; buffIndex < mDisplayWidth; ++buffIndex) {
      outputBuffer[buffIndex] = (buffIndex < textValue.length()) ? textValue[buffIndex] : ' ';
    }

    outputBuffer[mDisplayWidth] = '\0';

    mpLcdDisplayDevice->setCursor(0, lineIndex);
    mpLcdDisplayDevice->print(outputBuffer);

    delete[] outputBuffer;
}

void LcdDisplay::writeLine(int lineIndex, char* textValue) {
    if (textValue == nullptr) {
        throw std::runtime_error("LcdDisplay: Text value pointer is null");
    }

    return writeLine(lineIndex, String(textValue));
}
