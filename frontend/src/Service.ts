class Singleton {
    // Hold a reference to the single created instance
    // of the Singleton, initially set to null.
    private static instance: Singleton | null = null;

    // Make the constructor private to block instantiation
    // outside of the class.
    private constructor() {
        // initialization code
    }

    // Provide a static method that allows access
    // to the single instance.
    public static getInstance(): Singleton {
        // Check if an instance already exists.
        // If not, create one.
        if (this.instance === null) {
            this.instance = new Singleton();
        }
        // Return the instance.
        return this.instance;
    }

    // Example method to show functionality.
    public someMethod() {
        return "Doing something...";
    }
}
