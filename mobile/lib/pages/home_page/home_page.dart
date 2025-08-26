import 'package:flutter/material.dart';
import 'package:mobile/themes/home_page/colors.dart';

class HomePage extends StatefulWidget {
  const HomePage({super.key});

  @override
  State<HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> {
  //0 -> tasks
  //1-> home
  //2->messages
  int _selectedIndex = 1; // HOME is selected by default

  void _onItemTapped(int index) {
    setState(() {
      _selectedIndex = index;
    });
  }

  @override
  Widget build(BuildContext context) {

    String getCurrentDay(int day) {
      String name;
      switch (day) {
        case 1:
          name = "Monday";
          break;
        case 2:
          name = "Tuesday";
          break;
        case 3:
          name = "Wednesday";
          break;
        case 4:
          name = "Thursday";
          break;
        case 5:
          name = "Friday";
          break;
        case 6:
          name = "Saturday";
          break;
        case 7:
          name = "Sunday";
          break;
        default:
          name = "Unknown";
      }
      return name;
    }

    final size = MediaQuery.of(context).size;

    return Scaffold(
      backgroundColor: HomeColors.backgroundColor,
      body: Padding(
        padding: const EdgeInsets.all(24.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [

            const SizedBox(height: 70),
            // Greeting
            RichText(
              text: TextSpan(
                children: [
                  TextSpan(
                    text: 'Hi! ',
                    style: TextStyle(
                      fontSize: 75,
                      fontWeight: FontWeight.w600,
                      color: HomeColors.greetingTextColor,
                    ),
                  ),
                  TextSpan(
                    text: 'MARK',
                    style: TextStyle(
                      fontSize: 48,
                      fontWeight: FontWeight.w400,
                      color: HomeColors.greetingTextColor,
                    ),
                  ),
                ],
              ),
            ),


            RichText(
              text: TextSpan(
                children: [
                  TextSpan(
                    text: "Today is ",
                    style: TextStyle(
                      fontSize: 32,
                      fontWeight: FontWeight.w600,
                      color: HomeColors.dateTextColor,
                    )
                  ),
                  TextSpan(
                    text:getCurrentDay(DateTime.now().weekday),
                    style: TextStyle(
                      fontSize: 32,
                      fontWeight: FontWeight.w600,
                      color: HomeColors.dateTextColor,
                    )
                  ),
                ]
              ),
            ),

            const SizedBox(height:64),

            // Tasks section
            RichText(
              text: TextSpan(
                style: TextStyle(
                  fontSize: 24,
                  fontWeight: FontWeight.w400,
                  color: Colors.grey.shade600,
                ),
                children: [
                  TextSpan(
                      text: 'Today you have to do:',
                      style: TextStyle(
                        fontSize: 45,
                        fontWeight: FontWeight.w300,
                        color: HomeColors.textColor,
                    ),
                  ),
                  TextSpan(
                    text: ' 9 ',
                    style: TextStyle(
                      fontSize: 64,
                      fontWeight: FontWeight.w500,
                      color: HomeColors.boldNumColor,
                    ),
                  ),
                  TextSpan(
                    text: 'tasks',
                    style: TextStyle(
                      fontSize: 45,
                      fontWeight: FontWeight.w300,
                      color: HomeColors.textColor,
                    ),
                  ),
                ],
              ),
            ),

            const SizedBox(height: 32),

            // Messages section
            RichText(
              text: TextSpan(
                children:  [
                  TextSpan(
                    text: 'You have',
                    style: TextStyle(
                      fontSize: 45,
                      fontWeight: FontWeight.w300,
                      color: HomeColors.textColor,
                    ),
                  ),
                  TextSpan(
                    text: ' 0 ',
                    style: TextStyle(
                      fontSize: 64,
                      fontWeight: FontWeight.w500,
                      color: HomeColors.boldNumColor,
                    ),
                  ),
                  TextSpan(
                    text: 'new messages',
                    style: TextStyle(
                      fontSize: 45,
                      fontWeight: FontWeight.w300,
                      color: HomeColors.textColor,
                    ),
                  ),
                ],
              ),
            ),

          ],
        ),
      ),
      bottomNavigationBar: NavigationBar(
        selectedIndex: _selectedIndex,
        onDestinationSelected: _onItemTapped,
        backgroundColor: Colors.white,
        indicatorColor: Colors.transparent,
        height: 70,
        destinations: [
          NavigationDestination(
            icon:  Text(
              'TASKS',
              style: TextStyle(
                fontSize: 20,
                fontWeight: FontWeight.w600,
                color: HomeColors.boldNumColor,
              ),
            ),
            label: '',
          ),
          NavigationDestination(
            icon: Text(
              'HOME',
              style: TextStyle(
                fontSize: 20,
                fontWeight: FontWeight.w600,
                color: Colors.black,
              ),
            ),
            label: '',
          ),
          NavigationDestination(
            icon:  Text(
              'SEND',
              style: TextStyle(
                fontSize: 20,
                fontWeight: FontWeight.w600,
                color: HomeColors.boldNumColor,
              ),
            ),
            label: '',
          ),
        ],
      ),
    );
  }
}