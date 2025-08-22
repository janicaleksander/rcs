import 'package:flutter/material.dart';
import 'colors.dart';

ThemeData lightTheme = ThemeData(
  brightness: Brightness.light,
  scaffoldBackgroundColor: LoginColors.bodyColor,
  hintColor: LoginColors.textColor,
  primaryColorLight: LoginColors.buttonBackgroundColor,
  textTheme: TextTheme(
    headlineMedium: TextStyle(
      color: Colors.black,
      fontSize: 40,
      fontWeight: FontWeight.bold,
    )
  ),
  buttonTheme: ButtonThemeData(
    textTheme: ButtonTextTheme.primary,
    buttonColor: Colors.black,
  )
);