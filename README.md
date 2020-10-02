# 93L56R-CLI

A commandline tool for reading and writing 93L56R serial EEPROMs of the Combination Meter, and ECM, or the IS24C01 EEPROM used in the BIU, using an arduino
running my [subaru immo sketch](https://github.com/rgeyer/sketch_subaru_immo).

# TODO
* Tests?
* TravisCI or github actions to automate binary creation and publication
* Verbose output only with a flag.


Got all of the understanding of the odometer right here.
https://www.rs25.com/forums/f105/t105267-diy-reprogram-odometer-your-swapped-dash.html
https://www.rs25.com/forums/f105/1668064-post3.html
https://www.rs25.com/forums/f105/3250739-post81.html

Right-Vertical EEPROM contains the odometer on E0 and F0

Odometer value is a 20bit unsigned int, which will overflow at 1048576. The largest
usable number is 999999, since the odometer only has 6 decimal places. Not sure
how it would react to a number larger than 999999.
