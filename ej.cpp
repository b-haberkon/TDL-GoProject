class Lector() {


    private: 

    ifstream archivo;
    Lista<Escritor*>* listaEscritores;
    Lista<Lectura*>* listaLecturas;
    string nombreLectrura;


    public:

    Lector::Lector() {
        archivo.open("string"; std::in::out);
        escritores = new Lista<Escritor*>;
        lecturas = new Lista<Lectura*>;
    }


    Lista<Lectura*>* Lector::leerLectura() {

        while (archivo.isOpen()) {
                getline << Nombre;
        ...
        ...
        Lectura* lectura = new Lectura(...)
            agregarLista(lectura)
        }


        return this.listaLecturas;
 
    }


    Lector::agregarLista(Lectura* lectura) {
        this.listaEscritores.add(lectura);
    }




}


int main() {
    Lector* lector = new Lector;
    Lista<Lectura*>* = lector.leerLectura();
}