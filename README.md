# Proyecto Koi

Éste es un proyecto en Go para la materia 75.24 Teoría de la Programación / 75.31 Teoría del Lenguaje / 95.07 Teoría del Lenguaje de Programaciónes.
Consiste en un sitio web como herramienta complementaria para el aprendizaje del idioma japonés.

Para demostrar alguas características del lenguaje, hemos desarrollado un pequeño CRUS de usuarios, y un minijuego que emplea fuertemente la concurrencia.

## Cómo ejecutarlo
### Iniciar el servidor
Para iniciar el servidor, ejecutar `go run main.go` desde la raíz del proyecto.
El mismo correrá por defecto en el puerto 8000, luego de mostrar un pequeño recuadro con información.
Se vuelca también por consola varios informes para mostrar el funcionamiento.

Hay un [**Swagger/OpenAPI**](/doc/swagger.json) disponible.

### Iniciar el cliente
La primera vez es necesario instalar las dependencias, mediante el comando `npm i`.
Luego puede iniciarse con el comando `npm start`, que abrirá el puerto 3000.

## Sitio web

### Acceder al sitio web

## Minijuego Memotest
El minijuego puede accederse desde la misma máquina, en: http://127.0.0.1:3000/ejmemo.html

Consiste en un pequeño juego de memoria ("memotest") en el que se deben seleccionar de a dos piezas, intentando que sean del mismo par.
Las dos principales diferencias con el tradicional son que:

a. Es multijugador.

b. Los pares no son de piezas idénticas, sino de significados idénticos. Están compuestos por una palabra en japonés y un dibujo que la representa.
Pero, en ambos casos, desde el subdirectorio `react-auth/`.

### Cómo jugarlo
La interfaz es muy sencilla y minimalista ya que el foco del desarrollo está en el backend con concurrencia.
No deja de ser interesante la posibilidad de rediseñarlo usando en parte Golang compilando a [WASM (WebAssembly)](https://webassembly.org)

Para empezar hay que crear un juego (presionando en el botón), o unirse a un juego en curso ingresando previamente su número.

![Un cuadro para el id, un botón Unirse, y un botón Nuevo](doc/img/memo_inicial.PNG?raw=true "Pantalla inicial")

Para jugar, se deben cliquear las piezas ocultas (bordes lisos y finos, sin nada en el interior).
Las piezas seleccionadas tendrán el borde intermitente; de color rojo si la seleccionó uno, o azul si fue otro jugador.

![Una grilla con el primera pieza seleccionada por otro jugador (azul) y la segunda por uno mismo (rojo)](doc/img/memo_2_jugadores.PNG?raw=true "Dos jugadores")

Cuando el par de piezas seleccionadas por un jugador coincide, el borde se vuelve liso pero más intenso, y con el mismo código de colores.
Luego de unos segundos, las piezas son removidas, lo que se representa con una línea intermitente de color gris.

![Una grilla con de 4 filas y 3 columnas, las piezas 1, 2, 8 y 10 han sido removidas, y las piezas 3 y 5 forman una pareja de una planta y la palabra japonesa respectivamente](doc/img/memo_match.PNG?raw=true "Una coincidencia")

Las piezas removidas suman puntos para los jugadores, que pueden verse en la lista inferior ordenada de mayor a menor. El jugador actual se muestra en negrita.

Al finalizar el juego, ambos jugadores sólo pueden ver sus puntajes.

![Status: Ended. Y lista de puntos: "1. Anon2 (12)." En negrita "2. Anon1 (0)"](/doc/img/memo_ended.PNG?raw=true "Fin de juego")

## Diagramas de secuencia
### Seleccionar una pieza un solo jugador
![Diagrama de secuencia](/diagrams/__WorkspaceFolder__/doc/uc_select/uc_select_ok.png?raw=true "Un jugador (ok)")

### Dos jugadores intentan seleccionar la misma pieza
![Diagrama de secuencia](/diagrams/__WorkspaceFolder__/doc/uc_select/uc_select_collide.png?raw=true "Dos jugadores, una pieza (colisión)")

## Aplicación
La aplicación emplea una base de datos sqlite (aunque puede fácilmente cambiarse por mysql, porque se usa [gorm.io](https://gorm.io)).
