import 'package:dartz/dartz.dart';
import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:mobile/api/api_service.dart';
import 'package:mobile/api/err_response.dart';
import 'package:mobile/storage/storage_service.dart';
import 'package:get_it/get_it.dart';
import 'package:mobile/pages/tasks_page/tasks_page.dart' show UserTask;

Future<Either<List<UserTask>, ErrResponse>> getUserTasks(BuildContext ctx) async {
  final StorageService _storage = GetIt.instance<StorageService>();
  String? deviceID = await _storage.getValue("deviceID");

  if (deviceID == null) {
    return Right(ErrResponse("Device Error", "Device ID not found"));
  }

  try {
    final response = await ApiClient().dio.get("/tasks/$deviceID");
    if (response.statusCode == 200) {
      final tasksJson = response.data['tasks'] as List;
      final tasks = tasksJson.map((t) => UserTask.fromJson(t as Map<String, dynamic>)).toList();
      return Left(tasks);
    } else {
      return Right(ErrResponse(response.data['title'] ?? "Error", response.data['message'] ?? "Unknown error"));
    }
  } on DioException catch (e) {
    if (e.response != null) {
      final data = e.response?.data;
      return Right(ErrResponse(
          data['title'] ?? "Network Error",
          data['message'] ?? "Request failed"
      ));
    } else {
      return Right(ErrResponse("Connection Error", e.message ?? "Unable to connect to server"));
    }
  } catch (e) {
    return Right(ErrResponse("Unexpected Error", "An unexpected error occurred: ${e.toString()}"));
  }
}

Future<String> getNumOfTasks(BuildContext ctx) async{
  final result = await getUserTasks(ctx);
  return result.fold(
      (ok){
        return ok.length.toString();
      },
      (error){
        return "Check it!";
      }
  );
}