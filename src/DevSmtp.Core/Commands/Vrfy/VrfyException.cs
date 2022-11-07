namespace DevSmtp.Core.Commands
{
    public class VrfyException : Exception
    {
        public VrfyException(string message)
            : base(message)
        {
        }

        public VrfyException(string message, Exception innerException)
            : base(message, innerException)
        {
        }
    }
}
