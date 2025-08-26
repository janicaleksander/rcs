import 'package:dartz/dartz.dart';
import 'package:dio/dio.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:mobile/api/api_service.dart';
import 'package:get_it/get_it.dart';
import 'package:mobile/api/token_service.dart';
import 'package:mobile/api/err_response.dart';


class AuthService {
  final TokenService _storage = GetIt.instance<TokenService>();
  Future<Either<bool,ErrResponse>> login(String email, password) async {
    try {
      final response = await ApiClient().dio.post(
        "/login",
        data: {
          "email":email,
          "password":password,
        },
      );
      final responseCode = response.statusCode;
      switch (responseCode){
        case 200:
          final token = response.data['accessToken'];
          await _storage.saveToken(token);
          return Left(true);
        default:
          return Right(ErrResponse(response.data['title'], response.data['message']));
      }// this switch works onyl in range 200 - 299
    } on DioException catch (e) {
      if (e.response != null) {
        final data = e.response?.data;
        return Right(ErrResponse(data['title'], data['message']));

      } else {
        return Right(ErrResponse("Unexpected error!", e.message ?? "Unknown error"));
      }
    }

  }

  Future<void> logout() async {
    await _storage.deleteToken();
  }

}