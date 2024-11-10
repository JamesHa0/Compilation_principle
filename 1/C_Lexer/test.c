#include <stdio.h>

int addNumbers(int num1, int num2);

int main()
{
    int result;
    int num1 = 10;
    int num2 = 20;

    result = addNumbers(num1, num2);

    printf("两数相加的结果是：%d\n", result);

    return 0;
}

int addNumbers(int num1, int num2)
{
    return num1 + num2;
}