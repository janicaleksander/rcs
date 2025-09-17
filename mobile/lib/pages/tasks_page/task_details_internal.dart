import 'package:dartz/dartz.dart';
import 'package:flutter/material.dart';
import 'package:mobile/api/api_service.dart';
import 'package:mobile/api/err_response.dart';
import 'package:mobile/storage/storage_service.dart';
import 'package:get_it/get_it.dart';

//set at current
Future<Either<void,ErrResponse>> updateCurrentTask(BuildContext ctx,String taskID) async {
  final StorageService _storage = GetIt.instance<StorageService>();
  String? userID = await _storage.getValue("userID");
  String? deviceID = await _storage.getValue("deviceID");
  if (userID == null || deviceID == null) {
    return Right(ErrResponse("Error", "Empty params"));
  }

  final response = await ApiClient().dio.post(
      "/task/current/$userID?deviceID=$deviceID",
      data: {
        "taskID": taskID,
      }
  );
  final responseCode = response.statusCode;
  if (responseCode != 200) {
    final data = response.data;
    return Right(ErrResponse(
        data['title'] ?? "Network Error",
        data['message'] ?? "Request failed"
    ));
  }
  return Left(null);
}
// done task

Future<Either<void,ErrResponse>> deleteTask(BuildContext ctx,String taskID) async{
  final StorageService _storage  = GetIt.instance<StorageService>();
  String? deviceID = await _storage.getValue("deviceID");
  if (deviceID == null){
    return Right(ErrResponse("Error", "Empty param"));
  }
  final response = await ApiClient().dio.delete(
    "/task/$taskID?deviceID=$deviceID"
  );
  final responseCode = response.statusCode;
  if(responseCode != 200){
    final data = response.data;
    return Right(ErrResponse(data['title']??"Network error",data['message'] ?? "Request failed"));
  }
  return Left(null);
}