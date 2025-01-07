def caculatot():
    print("welcome to the python calculator!")
    print("select an opreation to perform")
    print("1. Addition")
    print("2.subtraction")
    print("3.multiplication")
    print("4.division")
    print("5.exit")
    while True:
        choice = input("\nEnter the number of your choice(1-5):")
        if choice == "5":
            print("exiting the calculator.goodbye!")
            break
        if choice in ["1,2,3,4"]:
            try:
                num1 = float(input("enter the first number:"))
                num2 = float(input("enter the second number:"))
                if choice == "1":
                    result = num1 + num2
                    print(f"the result of addition is: {result}")
                elif choice == "2":
                    result = num1 - num2
                    print(f"the result of subtraction is: {result}")
                elif choice == "3":
                    result = num1 * num2
                    print(f"the result of multiplication is:{result}")
                elif choice == "4":
                    num2 != 0
                    result = num1 / num2
                    print(f"the result of division is:{result}")
                else:
                    print("error: division by zero is not is allowed.")
            except ValueError:
                print("invalid input! please enter numeric values.")
                print("invalid choice! please select a valid option.")


exit()
from kivymd.app import MDApp
import webview


class HybridApp(MDApp):
    def build(self):
        self.theme_cls.primary_palette = "Blue"
        return None  # Add a layout for UI if needed

    def on_start(self):
        webview.create_window("My Website", "http://localhost:1481/home/")
        webview.start()


if __name__ == "__main__":
    HybridApp().run()
