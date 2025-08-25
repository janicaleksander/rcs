import 'package:dio/dio.dart';
import 'package:flutter_dotenv/flutter_dotenv.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:get_it/get_it.dart';
import 'package:mobile/api/token_service.dart';

class ApiClient {
  final TokenService _storage = GetIt.instance<TokenService>();
  late final Dio dio;

  ApiClient() {
    dio = Dio(
      BaseOptions(
        baseUrl: "http://10.0.2.2:8978",

        // baseUrl: "http://127.0.0.1:8978",
       // baseUrl: dotenv.env['API_URL'] ?? "localhost",
        connectTimeout: const Duration(seconds: 10),
        receiveTimeout: const Duration(seconds: 10),
      ),
    )..interceptors.add(
      InterceptorsWrapper(
        onRequest: (options, handler) async {
          final token = await _storage.getToken();
          if (token != null) {
            options.headers['Authorization'] = 'Bearer $token';
          }
          return handler.next(options);
        },
        onError: (DioException e,handler) async{
          if(e.response?.statusCode == 401){
            //TODO refresg token
          }
          return handler.next(e);
        }
      ),
    );
  }
}
