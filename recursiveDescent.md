# Recursive Descent 

ex. "a(cat|cow)*"

grammar for regex:
E  -->  E -> E "|" E 
        E -> EE //use + in parse tree 
        E -> E "*" 
        E -> "(" E ")"
        v //any constant
