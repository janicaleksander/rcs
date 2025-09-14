import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:mobile/pages/home_page/home_page.dart';
import 'package:mobile/pages/login_page/login_page.dart';
import 'package:mobile/pages/tasks_page/tasks_page.dart';

final GoRouter router = GoRouter(
  initialLocation: '/login',
    routes: <RouteBase>[
    GoRoute(
      path: '/login',
      builder: (BuildContext context,GoRouterState state){
        return const LoginPage();
      },
    ),
    GoRoute(  
        path: '/home',
        builder: (BuildContext context,GoRouterState state){
          return const HomePage();
        },
    ),
    GoRoute(
        path: '/tasks',
        builder: (BuildContext context,GoRouterState state){
          return const TasksPage();
        },
    ),
  ]
);