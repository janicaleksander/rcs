import 'package:get_it/get_it.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:mobile/api/token_service.dart';

final getIt =  GetIt.instance;
void setupLocator(){
  getIt.registerLazySingleton<TokenService>(
      () =>  TokenService(),
  );
}