import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:mobile/pages/tasks_page/tasks_internal.dart';
import 'package:awesome_snackbar_content/awesome_snackbar_content.dart';

class UserTask {
  final String id;
  final String name;
  final String description;
  final int state;
  final DateTime deadline;
  final IconData icon;

  UserTask({
    required this.id,
    required this.name,
    required this.description,
    required this.state,
    required this.deadline,
    this.icon = Icons.favorite,
  });

  factory UserTask.fromJson(Map<String, dynamic> json) {
    final deadlineJson = json['deadline'];
    final deadline = DateTime.fromMillisecondsSinceEpoch(
      (deadlineJson['seconds'] * 1000) + (deadlineJson['nanos'] ~/ 1000000),
      isUtc: true,
    );

    return UserTask(
      id: json['id'].toString(),
      name: json['name'],
      description: json['description'],
      state: int.tryParse(json['state'].toString()) ?? 0,
      deadline: deadline,
    );
  }
}

SnackBar createError(String title, String err) {
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
    ),
  );
}

class TasksPage extends StatefulWidget {
  const TasksPage({super.key});

  @override
  State<TasksPage> createState() => _TasksPageState();
}

class _TasksPageState extends State<TasksPage> {
  Future<List<UserTask>>? _tasksFuture;

  Future<List<UserTask>> fetchTasksFromDatabase(BuildContext ctx) async {
    final result = await getUserTasks(ctx);

    return await result.fold(
          (tasks) async {
        return tasks;
      },
          (err) async {
        SnackBar error = SnackBar(content: Text(err.message));
        ScaffoldMessenger.of(ctx)
          ..hideCurrentSnackBar()
          ..showSnackBar(error);
        return [];
      },
    );
  }

  @override
  void initState() {
    super.initState();
    _tasksFuture = fetchTasksFromDatabase(context);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => context.pop(),
        ),
        title: const Text(
          "Your tasks!",
          style: TextStyle(fontWeight: FontWeight.bold),
        ),
        backgroundColor: Colors.white,
        elevation: 0,
      ),
      body: FutureBuilder<List<UserTask>>(
        future: _tasksFuture,
        builder: (context, snapshot) {
          if (snapshot.connectionState == ConnectionState.waiting) {
            return const Center(child: CircularProgressIndicator());
          }
          if (!snapshot.hasData || snapshot.data!.isEmpty) {
            return const Center(child: Text("No tasks found"));
          }

          final tasks = snapshot.data!;

          return ListView.builder(
            padding: const EdgeInsets.all(8),
            itemCount: tasks.length,
            itemBuilder: (context, index) {
              final task = tasks[index];
              return Padding(
                padding: const EdgeInsets.symmetric(vertical: 4),
                child: Card(
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(12),
                  ),
                  color: Colors.blue.withOpacity(0.1),
                  child: ListTile(
                    leading: CircleAvatar(
                      backgroundColor: Colors.white,
                      child: Icon(task.icon, color: Colors.blue, size: 20),
                    ),
                    title: Text(task.name,
                        style: const TextStyle(fontWeight: FontWeight.bold)),
                    subtitle: Text(
                      "${task.deadline.day}-${task.deadline.month}-${task.deadline.year} "
                          "${task.deadline.hour}:${task.deadline.minute.toString().padLeft(2, '0')}",
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                    trailing: const Icon(Icons.arrow_forward_ios, size: 14),
                    onTap: () {
                      context.push('/taskd', extra: task).then((_) {
                        setState(() {
                          _tasksFuture = fetchTasksFromDatabase(context);
                        });
                      });
                    },
                  ),
                ),
              );
            },
          );
        },
      ),
    );
  }
}
