import 'package:dartz/dartz.dart';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:mobile/api/auth_service.dart';
import 'package:mobile/themes/login_page/colors.dart';
import 'package:awesome_snackbar_content/awesome_snackbar_content.dart';
import 'package:http/http.dart' as http;
SnackBar createError(String title,err){
  return SnackBar(
      elevation: 0,
      behavior: SnackBarBehavior.floating,
      backgroundColor: Colors.transparent,
      content: SizedBox(
        height: 100,
        child: AwesomeSnackbarContent(
          title: title,
          message: err,
          contentType: ContentType.failure,
        ),
      )
  );
}



Future<void> onLoginPressed(BuildContext ctx,TextEditingController email, TextEditingController password) async{
  final emailText = email.text.trim();
  final passwordText = password.text.trim();
  final result = await AuthService().login(emailText, passwordText);
  if (!ctx.mounted) return;
  
  
  result.fold(
      (ok) {
        SnackBar error = createError("Login error","SUCCESS");
        ScaffoldMessenger.of(ctx)
          ..hideCurrentSnackBar()
          ..showSnackBar(error);
      },
      (error){
        SnackBar err = createError(error.title,error.message);
        ScaffoldMessenger.of(ctx)
          ..hideCurrentSnackBar()
          ..showSnackBar(err);
      },
  );
  

  /*
  result.fold(
      (response){
      ctx.go('/');//TODO
      },
      (errorMsg){
        SnackBar error = createError("Login error",errorMsg);
        ScaffoldMessenger.of(ctx)
          ..hideCurrentSnackBar()
          ..showSnackBar(error);
      },
  );

   */


}
