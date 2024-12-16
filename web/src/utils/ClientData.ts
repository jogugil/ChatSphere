export async function getClientInformation (): Promise<string> {
    // Obtener detalles del navegador
    const navegador = {
        userAgent: navigator.userAgent, // Utiliza userAgent en lugar de appName
        language: navigator.language,  
        vendor: navigator.vendor, // Añadir el proveedor del navegador
      };

    // Obtener la IP y el país usando una API de geolocalización
    const getLocation = async (): Promise<{ ip: string, pais: string }> => {
        try {
            const respuesta1 = await fetch('https://cors-anywhere.herokuapp.com/https://ipapi.co/json/');
            const respuesta2 = await fetch('https://ipapi.co/json/');
            const datos = await respuesta1.json();
            console.log('Datos obtenidos de ipapi.co:', datos); // Añadir log para verificar la respuesta
            return { ip: datos.ip, pais: datos.country_name };
        } catch (error) {
            console.log('Error al obtener la ubicación:', error);
            return { ip: 'Desconocida', pais: 'Desconocido' };
        }
    };

    // Obtener IP y país
    const { ip, pais } = await getLocation();

    // Crear un objeto con toda la información
    const informacionCliente = {
        navegador: navegador,
        ip: ip,
        pais: pais,
    };

    // Convertir la información a un formato de string JSON
    return JSON.stringify(informacionCliente, null, 2); // Con formato legible
}
