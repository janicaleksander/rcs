import 'package:get_it/get_it.dart';
import 'package:mobile/api/token_service.dart';
import 'package:mobile/location/location_service.dart';
import 'package:mobile/storage/storage_service.dart';

final getIt =  GetIt.instance;
void setupLocator(){
  getIt.registerLazySingleton<TokenService>(
      () =>  TokenService(),
  );
  getIt.registerLazySingleton<LocationService>(
      () =>  LocationService(),
  );
  getIt.registerLazySingleton<StorageService>(
      () => StorageService(),
  );
}