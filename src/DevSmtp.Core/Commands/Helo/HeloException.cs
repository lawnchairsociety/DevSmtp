namespace DevSmtp.Core.Commands
{
    public class HeloException : Exception
    {
        public HeloException(string message)
            : base(message)
        {
        }

        public HeloException(string message, Exception innerException)
            : base(message, innerException)
        {
        }
    }
}
