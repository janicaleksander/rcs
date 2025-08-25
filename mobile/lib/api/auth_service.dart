import 'package:dio/dio.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:mobile/api/api_service.dart';
import 'package:get_it/get_it.dart';
import 'package:mobile/api/token_service.dart';



class AuthService {
  final TokenService _storage = GetIt.instance<TokenService>();
  Future<bool> login(String email, password) async {
    try {
      final response = await ApiClient().dio.post(
        "/login",
        data: {
          "email":email,
          "password":password,
        },
      );

      final token = response.data['accessToken'];
      if(token != null) {
        await _storage.saveToken(token);
        return true;
      }
      return false;

    } on DioException catch (e){
      print(e);
      return false;
      //error
    }
  }

  Future<void> logout() async {
    await _storage.deleteToken();
  }

}