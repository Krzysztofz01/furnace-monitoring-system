#include <OneWire.h>
#include <LiquidCrystal.h>

OneWire temperature_sensor_one(D3);
OneWire temperature_sensor_two(D4);
LiquidCrystal lcd_screen(D1, D2, D5, D6, D7, D8);

const int LCD_WIDTH = 16;
const int LCD_HEIGHT = 2;

class TemperatureResult
{
public:
  static TemperatureResult Success(float temperature) {
    TemperatureResult result;
    result.mFailureMessage = "";
    result.mIsSuccess = true;
    result.mTemperature = temperature;
    
    return result;
  }

  static TemperatureResult Failure(String message) {
    TemperatureResult result;
    result.mFailureMessage = message;
    result.mIsSuccess = false;
    result.mTemperature = 0.0f;

    return result;
  }

  bool is_success() {
    return mIsSuccess;
  }

  float get_temperature() {
    return mTemperature;
  }

  String get_failure_message() {
    return mFailureMessage;
  }

  ~TemperatureResult() { }
private:
  TemperatureResult() { }

  String mFailureMessage;
  bool mIsSuccess;
  float mTemperature;
};

// TODO: Preserve some kind of identifier on failure messages
TemperatureResult get_sensor_temperature(OneWire* onewire_device) {
  const int address_buffer_size = 8;
  byte address_buffer[address_buffer_size];
  
  if (!onewire_device->search(address_buffer)) {
    onewire_device->reset_search();
    delay(250);
    return TemperatureResult::Failure("No addresses found to read sensor data.");
  }

  if (OneWire::crc8(address_buffer, 7) != address_buffer[7]) {
    onewire_device->reset_search();
    return TemperatureResult::Failure("Verification of address CRC checksum failed.");
  }

  // NOTE: Identification of the sensor circuit based on the first communication byte
  byte ic_identifier;
  switch (address_buffer[0]) {
    // NOTE: DS118S20 or old DS1820
    case 0x10:
      ic_identifier = 1;
      break;

    // NOTE: DS18B20
    case 0x28:
      ic_identifier = 0;
      break;

    // NOTE: DS1822
    case 0x22:
      ic_identifier = 0;
      break;

    // NOTE: Undefined sensor circuit
    default:
      return TemperatureResult::Failure("The sensor circuit is not supported");
  }

  onewire_device->reset();
  onewire_device->select(address_buffer);
  onewire_device->write(0x44, 1);

  // TODO: This value can be fine-tuned down to 750ms?
  delay(1000);

  byte present = 0;
  present = onewire_device->reset();
  onewire_device->select(address_buffer);
  onewire_device->write(0xBE);

  const int value_buffer_size = 12;
  byte value_buffer[value_buffer_size];

  for (int i = 0; i < value_buffer_size - 3; ++i) {
    value_buffer[i] = onewire_device->read();
  }

  onewire_device->read();
  onewire_device->reset_search();

  // NOTE: Parsing the actual temperature for the data buffer. The result data size
  //       is always a 16 bit signed integer, and the type MUST be explicit in order
  //       to make the firmware work also for 32-bit processor's.
  int16_t raw_data = (value_buffer[1] << 8) | value_buffer[0];

  if (ic_identifier) {
    raw_data = raw_data << 3;

    if (value_buffer[7] == 0x10) {
      raw_data = (raw_data & 0xFFF0) + value_buffer_size - value_buffer[6];
    }
  } else {
    byte cfg = (value_buffer[4] & 0x60);

    if (cfg == 0x00) {
      // NOTE: 9bit resolution 94ms
      raw_data = raw_data & ~7;
    } else if (cfg == 0x20) {
      // NOTE: 10bit resolution 188ms
      raw_data = raw_data & ~3;
    } else if (cfg == 0x40) {
      // NOTE: 11bit resolution 375ms
      raw_data = raw_data & ~1;
    }

    // NOTE: Default resolution 750ms
  }

  float celsius_temperature = (float)raw_data / 16.0;
  return TemperatureResult::Success(celsius_temperature);
}

void setup(void) {
  Serial.begin(9600);
  lcd_screen.begin(LCD_WIDTH, LCD_HEIGHT);
}
 
void loop(void) {
  lcd_screen.setCursor(0, 0);
  auto result_ts_one = get_sensor_temperature(&temperature_sensor_one);
  if (result_ts_one.is_success()) {
    char print_buffer[LCD_WIDTH];
    snprintf(print_buffer, LCD_WIDTH, "Temp one: %.4f", result_ts_one.get_temperature());

    Serial.println(print_buffer);
    lcd_screen.print(print_buffer);
  }

  lcd_screen.setCursor(0, 1);
  auto result_ts_two = get_sensor_temperature(&temperature_sensor_two);
  if (result_ts_two.is_success()) {
    char print_buffer[LCD_WIDTH];
    snprintf(print_buffer, LCD_WIDTH, "Temp two: %.4f", result_ts_two.get_temperature());

    Serial.println(print_buffer);
    lcd_screen.print(print_buffer);
  }

  Serial.print("\n");
}