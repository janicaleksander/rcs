import 'package:flutter_secure_storage/flutter_secure_storage.dart';


class StorageService {
  final _storage  = const FlutterSecureStorage();
  Future<void> save(String key, String value) async{
    await _storage.write(key: key,value: value);
  }

  Future<String?> getValue(String key) async {
    return await _storage.read(key: key);
  }

  Future<void> deleteKey(String key) async {
    await _storage.delete(key: key);
  }
}