import 'package:flutter/material.dart';
import 'package:mobile/themes/login_page/login_theme.dart';
import 'package:mobile/pages/login_page/login_page.dart';
void main(){
  runApp(MaterialApp(
    home: LoginPage(),
    theme: LoginTheme.light,
    darkTheme: LoginTheme.dark,
    themeMode: ThemeMode.system,
    debugShowCheckedModeBanner: false,
  ));
}