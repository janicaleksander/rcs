import 'package:flutter/material.dart';
import 'package:mobile/router/router.dart';
import 'package:mobile/themes/login_page/login_theme.dart';
import 'package:mobile/pages/login_page/login_page.dart';


void main() => runApp(const MainApp());

class MainApp extends StatelessWidget {
  const MainApp({super.key});

  @override
  Widget build(BuildContext context) => MaterialApp.router(
    routerConfig: router,
    theme: LoginTheme.light,
    themeMode: ThemeMode.system,
    debugShowCheckedModeBanner: false,
  );

}
