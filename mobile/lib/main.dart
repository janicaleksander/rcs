import 'package:flutter/material.dart';
import 'package:mobile/api/service_locator.dart';
import 'package:mobile/router/router.dart';
import 'package:mobile/themes/login_page/login_theme.dart';
import 'package:mobile/pages/login_page/login_page.dart';
import 'package:flutter_dotenv/flutter_dotenv.dart';

void main() async{
  WidgetsFlutterBinding.ensureInitialized();
 /*
  try{
    await dotenv.load(fileName: ".env");
  }catch (e){
    throw Exception("Error with loading .env file!");
  }

  */
  setupLocator();

  runApp(const MainApp());
}
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
