namespace DevSmtp.Core.Commands
{
    public class QuitException : Exception
    {
        public QuitException(string message)
            : base(message)
        {
        }

        public QuitException(string message, Exception innerException)
            : base(message, innerException)
        {
        }
    }
}
