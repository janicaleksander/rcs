import 'package:flutter/material.dart';
import 'colors.dart';

ThemeData darkTheme = ThemeData(
    brightness: Brightness.dark,
    scaffoldBackgroundColor: LoginColors.bodyColorDark,
    hintColor: LoginColors.textColorDark,
    primaryColorLight: LoginColors.buttonBackgroundColorDark,
    textTheme: TextTheme(
        headlineMedium: TextStyle(
          color: Colors.white,
          fontSize: 40,
          fontWeight: FontWeight.bold,
        )
    ),
    buttonTheme: ButtonThemeData(
      textTheme: ButtonTextTheme.primary,
      buttonColor: Colors.white,
    )
);