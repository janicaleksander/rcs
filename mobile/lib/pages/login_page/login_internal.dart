import 'package:dartz/dartz.dart';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
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

Future<Either<http.Response,String>> loginUser(String email,password) async{
  try {
    return Right("xd");
  } catch (e){
    return Right(e.toString());
  }
}

Future<void> onLoginPressed(BuildContext ctx,TextEditingController email, TextEditingController password) async{
  final emailText = email.text.trim();
  final passwordText = password.text.trim();
  final result = await loginUser(emailText, passwordText);
  if (!ctx.mounted) return;
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


}
