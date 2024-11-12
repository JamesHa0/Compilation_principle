#include <stdio.h>

int addNumbers(int num1, int num2);

int main()
{
    // Single-line comments
    int result;
    int 2num1 = 10;
    int num2 = 20;
    double num3 = 5.00;
    /* Multi-line comments */
    result = addNumbers(num1, num2);

    printf("两数相加的结果是：%d\n", result);

    return 0;
}

int addNumbers(int num1, int num2)
{
    return num1 + num2;
}