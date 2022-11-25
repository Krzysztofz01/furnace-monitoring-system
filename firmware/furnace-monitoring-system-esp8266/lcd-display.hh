#ifndef LCD_DISPLAY_HH
#define LCD_DISPLAY_HH

#include <LiquidCrystal.h>
#include <stdexcept>

class LcdDisplay {
public:
    void clear();
    void writeLine(int lineIndex, String textValue);
    void writeLine(int lineIndex, char* textValue);

    LcdDisplay(const int displayWidth, const int displayHeight, const int pinRegisterSelect, const int pinEnable, const int pinDataBus1, const int pinDataBus2, const int pinDataBus3, const int pinDataBus4);
    ~LcdDisplay();
private:
    LcdDisplay();

    LiquidCrystal* mpLcdDisplayDevice;
    int mDisplayWidth;
    int mDisplayHeight;
};

#endif