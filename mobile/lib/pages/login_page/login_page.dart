import 'package:flutter/material.dart';
import 'package:mobile/themes/login_page/colors.dart';
import 'package:awesome_snackbar_content/awesome_snackbar_content.dart';

class LoginPage extends StatefulWidget {
  const LoginPage({super.key});

  @override
  State<LoginPage> createState() => _LoginPageState();
}

class _LoginPageState extends State<LoginPage> {
  final _emailController = TextEditingController();
  final _passwordController = TextEditingController();

  @override
  void dispose() {
    _emailController.dispose();
    _passwordController.dispose();
    super.dispose();
  }
  void _onLoginPressed() {
    final email = _emailController.text.trim();
    final password = _passwordController.text.trim();

    debugPrint("Email: $email, Password: $password");

    final snackBar = SnackBar(
      elevation: 0,
      behavior: SnackBarBehavior.floating,
      backgroundColor: Colors.transparent,
      content: SizedBox(
        height: 100,
        child: AwesomeSnackbarContent(
          title: 'Login Failed!',
          message: 'Invalid email or password. Please try again.',
          contentType: ContentType.failure,
        ),
      )


    );

    ScaffoldMessenger.of(context)
      ..hideCurrentSnackBar()
      ..showSnackBar(snackBar);
  }

  @override
  Widget build(BuildContext context) {
    final size = MediaQuery.of(context).size;

    return Scaffold(
      backgroundColor: LoginColors.backgroundColor,
      body: SafeArea(
        child: SingleChildScrollView(
          padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 60),
          child: SizedBox(
            width: size.width,
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                // Welcome Text
                RichText(
                  text: TextSpan(
                    children: [
                      TextSpan(
                        text: "WELCOME!\n",
                        style: Theme.of(context).textTheme.headlineLarge?.copyWith(
                          color: LoginColors.textColor,
                          fontWeight: FontWeight.bold,
                          fontSize: 64,
                          height: 1,
                        ),
                      ),
                      TextSpan(
                        text: "LET’S DO SOME\n",
                        style: Theme.of(context).textTheme.headlineSmall?.copyWith(
                          color: LoginColors.textColor,
                          fontWeight: FontWeight.w500,
                          fontSize: 40,
                          height: 1,
                        ),
                      ),
                      TextSpan(
                        text: "WORK",
                        style: Theme.of(context).textTheme.headlineSmall?.copyWith(
                          color: LoginColors.textHighlightColor,
                          fontWeight: FontWeight.w600,
                          fontSize: 72,
                          height: 0.8
                        ),
                      ),
                    ],
                  ),
                ),

                const SizedBox(height: 70),

                // Email Field
                TextField(
                  controller: _emailController,
                  keyboardType: TextInputType.emailAddress,
                  decoration: InputDecoration(
                    hintText: "Email",
                    prefixIcon: const Icon(Icons.email_outlined),
                    border: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(12),
                      borderSide: BorderSide(
                          color: LoginColors.textInputLineColor,
                          width: 3.0,
                      ), // kolor obramówki
                    ),
                    enabledBorder: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(12),
                      borderSide: BorderSide(
                          color: LoginColors.textInputLineColor,
                          width: 3.0,
                      ), // kolor gdy pole nieaktywne
                    ),
                    focusedBorder: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(12),
                      borderSide: BorderSide(
                          color: LoginColors.textInputLineColor,
                          width: 3.0,
                      ), // kolor gdy pole aktywne
                    ),
                  ),
                ),
                const SizedBox(height: 16),

                // Password Field
                TextField(
                  controller: _passwordController,
                  obscureText: true,
                  decoration: InputDecoration(
                    hintText: "Password",
                    prefixIcon: const Icon(Icons.lock_outline),
                    border: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(12),
                      borderSide: BorderSide(
                        color: LoginColors.textInputLineColor,
                        width: 3.0,
                      ), // kolor obramówki
                    ),
                    enabledBorder: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(12),
                      borderSide: BorderSide(
                        color: LoginColors.textInputLineColor,
                        width: 3.0,
                      ), // kolor gdy pole nieaktywne
                    ),
                    focusedBorder: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(12),
                      borderSide: BorderSide(
                        color: LoginColors.textInputLineColor,
                        width: 3.0,
                      ), // kolor gdy pole aktywne
                    ),
                  ),
                ),

                const SizedBox(height: 70),

                // Login Button
                SizedBox(
                  width: double.infinity,
                  height: 55,
                  child: ElevatedButton(
                    onPressed: _onLoginPressed,
                    style: ElevatedButton.styleFrom(
                      backgroundColor: LoginColors.buttonBackgroundColor,
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(12),
                      ),
                    ),
                    child: const Text(
                      "LOGIN",
                      style: TextStyle(
                        fontWeight: FontWeight.bold,
                        fontSize: 20,
                        color: Colors.white,
                      ),
                    ),
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
