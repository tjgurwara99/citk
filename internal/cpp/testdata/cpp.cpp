#include <iostream>

namespace something {
	int this_is_my_func() {
    	return 0;
    }

    class Car {
  		public:
	    int speed(int maxSpeed);

        void Something() {
            std::cout << "hello";
        };
	};

    int Car::speed(int maxSpeed) {
  		return maxSpeed;
	}
}
