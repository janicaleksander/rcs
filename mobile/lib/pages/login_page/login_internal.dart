import 'package:dartz/dartz.dart';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:mobile/api/api_service.dart';
import 'package:mobile/api/auth_service.dart';
import 'package:mobile/location/location_service.dart';
import 'package:mobile/storage/storage_service.dart';
import 'package:mobile/themes/login_page/colors.dart';
import 'package:awesome_snackbar_content/awesome_snackbar_content.dart';
import 'package:http/http.dart' as http;
import 'package:get_it/get_it.dart';

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
  final LocationService _locationService = GetIt.instance<LocationService>();
  final StorageService _storageService = GetIt.instance<StorageService>();
  if (!ctx.mounted) return;

  result.fold(
      (ok) async {
        bool isPermission = await _locationService.isPermission();
        if(!isPermission){
          SnackBar err = createError("Location error","Turn on location service!");
          ScaffoldMessenger.of(ctx)
            ..hideCurrentSnackBar()
            ..showSnackBar(err);
          return;
        }
       bool isWorking  = await _locationService.enableBackgroundMode();
        if(!isWorking){return;}
        _locationService.registerHandler(() async {
          final response = await ApiClient().dio.post(
            "/location",
            data: {
              "location": {
                "latitude": (await _locationService.getLocationData())?.latitude,
                "longitude": (await _locationService.getLocationData())?.longitude,
              },
              "deviceID": (await _storageService.getValue('deviceID')),
            },
          );

          final responseCode = response.statusCode;
          if (responseCode != 200) {
            print("ERROR");
            return;
          }
        });

        ctx.go('/home');
      },
      (error){
        SnackBar err = createError(error.title,error.message);
        ScaffoldMessenger.of(ctx)
          ..hideCurrentSnackBar()
          ..showSnackBar(err);
      },
  );
  



}
