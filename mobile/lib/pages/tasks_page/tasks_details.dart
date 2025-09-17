import 'package:dartz/dartz.dart' hide State;
import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:mobile/api/api_service.dart';
import 'package:mobile/api/err_response.dart';
import 'package:mobile/pages/tasks_page/task_details_internal.dart';
import 'package:mobile/storage/storage_service.dart';
import 'package:get_it/get_it.dart';
import 'package:mobile/pages/tasks_page/tasks_page.dart' show UserTask;
import 'package:awesome_snackbar_content/awesome_snackbar_content.dart';

SnackBar createInfo(String title,desc){
  return SnackBar(
      elevation: 0,
      behavior: SnackBarBehavior.floating,
      backgroundColor: Colors.transparent,
      content: SizedBox(
        height: 100,
        child: AwesomeSnackbarContent(
          title: title,
          message: desc,
          contentType: ContentType.success,
        ),
      )
  );
}

class TaskDetail extends StatefulWidget {
  final UserTask? task;

  const TaskDetail({super.key, this.task});

  @override
  State<TaskDetail> createState() => _TaskDetailState();
}

class _TaskDetailState extends State<TaskDetail> {
  @override
  Widget build(BuildContext context) {
    if (widget.task == null) {
      return Scaffold(
        appBar: AppBar(
          title: const Text("Task Details"),
          leading: IconButton(
            icon: const Icon(Icons.arrow_back),
            onPressed: () => context.pop(),
          ),
        ),
        body: const Center(
          child: Text("No task data available"),
        ),
      );
    }

    final task = widget.task!;

    return Scaffold(
      appBar: AppBar(
        title: Text(task.name),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => context.pop(),
        ),
      ),
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              task.name,
              style: const TextStyle(
                  fontSize: 24,
                  fontWeight: FontWeight.bold
              ),
            ),
            const SizedBox(height: 16),

            Text(
              "Description:",
              style: TextStyle(
                fontSize: 16,
                fontWeight: FontWeight.w600,
                color: Colors.grey[700],
              ),
            ),
            const SizedBox(height: 8),
            Text(
              task.description,
              style: const TextStyle(fontSize: 16),
            ),
            const SizedBox(height: 24),

            Text(
              "Deadline:",
              style: TextStyle(
                fontSize: 16,
                fontWeight: FontWeight.w600,
                color: Colors.grey[700],
              ),
            ),
            const SizedBox(height: 8),
            Text(
              "${task.deadline.day}/${task.deadline.month}/${task.deadline.year} ${task.deadline.hour}:${task.deadline.minute.toString().padLeft(2, '0')}",
              style: const TextStyle(fontSize: 16),
            ),
            const SizedBox(height: 24),

            // Status
            ElevatedButton(
              onPressed: () async {
                final resp = await updateCurrentTask(context, task.id);
                resp.fold(
                      (_) {
                    final info = createInfo("Success", "Task set as current!");
                    ScaffoldMessenger.of(context)
                      ..hideCurrentSnackBar()
                      ..showSnackBar(info);
                  },
                      (err) {
                    final info = createInfo(err.title, err.message);
                    ScaffoldMessenger.of(context)
                      ..hideCurrentSnackBar()
                      ..showSnackBar(info);
                  },
                );
              },
              child: const Text("Set as current!"),
            ),
            ElevatedButton(
              onPressed: () async {
                final resp = await deleteTask(context, task.id);
                resp.fold(
                      (_) {
                    final info = createInfo("Success", "Task finished!");
                    ScaffoldMessenger.of(context)
                      ..hideCurrentSnackBar()
                      ..showSnackBar(info);
                    context.pop();
                  },
                      (err) {
                    final info = createInfo(err.title, err.message);
                    ScaffoldMessenger.of(context)
                      ..hideCurrentSnackBar()
                      ..showSnackBar(info);
                  },
                );
              },
              child: const Text("Finish task!"),
            ),

            const SizedBox(height: 8),
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
              decoration: BoxDecoration(
                color: task.state == 0 ? Colors.orange[100] : Colors.green[100],
                borderRadius: BorderRadius.circular(8),
              ),
            ),
          ],
        ),
      ),
    );
  }
}